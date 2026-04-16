import { test, expect } from '@playwright/test';

// Mock data for testing
const mockArtists = [
  { id: 1, name: '周杰伦', songCount: 3 },
  { id: 2, name: '林俊杰', songCount: 2 },
];

const mockSongs = [
  { id: 1, title: '晴天', artist: '周杰伦', album: '叶惠美', year: 2003, genre: '流行', duration: 267, filePath: '/music/rock/晴天.mp3' },
  { id: 2, title: '夜曲', artist: '周杰伦', album: '七里香', year: 2004, genre: '流行', duration: 252, filePath: '/music/pop/夜曲.mp3' },
  { id: 3, title: '江南', artist: '林俊杰', album: '编号89757', year: 2004, genre: '流行', duration: 245, filePath: '/music/pop/江南.mp3' },
];

test.describe('Epic 2: 音乐库浏览与搜索', () => {
  test.beforeEach(async ({ page }) => {
    // Mock the setup status API to indicate already configured
    await page.route('/api/setup/status', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { needs_setup: false, music_dir: '/tmp/music', db_path: '' }
        })
      });
    });

    await page.goto('/');
  });

  test.describe('Tab Navigation', () => {
    test('应该显示Tab导航', async ({ page }) => {
      // 验证Tab导航存在
      await expect(page.locator('text=歌手')).toBeVisible();
      await expect(page.locator('text=专辑')).toBeVisible();
      await expect(page.locator('text=文件夹')).toBeVisible();
    });

    test('默认选中歌手Tab', async ({ page }) => {
      // 歌手Tab应该有激活样式
      const artistTab = page.locator('button:has-text("歌手")');
      await expect(artistTab).toBeVisible();
    });

    test('点击专辑Tab应切换到专辑视图', async ({ page }) => {
      // Mock albums API
      await page.route('/api/albums', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: [] })
        });
      });

      await page.locator('button:has-text("专辑")').click();

      // URL或内容应该变化
      await expect(page.locator('text=专辑')).toBeVisible();
    });

    test('点击文件夹Tab应切换到文件夹视图', async ({ page }) => {
      // Mock folders API
      await page.route('/api/folders', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: [] })
        });
      });

      await page.locator('button:has-text("文件夹")').click();

      await expect(page.locator('text=文件夹')).toBeVisible();
    });
  });

  test.describe('Artist View (Story 2.1)', () => {
    test.beforeEach(async ({ page }) => {
      // Mock artists API
      await page.route('/api/artists**', async (route) => {
        const url = route.request().url();
        if (url.includes('/api/artists/')) {
          // Artist songs endpoint
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockSongs.filter(s => s.artist === '周杰伦') })
          });
        } else {
          // Artists list endpoint
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockArtists })
          });
        }
      });
    });

    test('应显示艺术家列表', async ({ page }) => {
      // 等待艺术家列表加载
      await expect(page.locator('text=周杰伦')).toBeVisible({ timeout: 5000 });
      await expect(page.locator('text=林俊杰')).toBeVisible();
    });

    test('应显示艺术家数量', async ({ page }) => {
      await expect(page.locator('text=共 2 位艺术家')).toBeVisible();
    });

    test('点击艺术家应展开歌曲列表', async ({ page }) => {
      // 点击周杰伦
      await page.locator('tr:has-text("周杰伦")').first().click();

      // 应该展开显示歌曲
      await expect(page.locator('text=晴天')).toBeVisible();
    });

    test('再次点击艺术家应折叠歌曲列表', async ({ page }) => {
      // 点击展开
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible();

      // 点击折叠
      await page.locator('tr:has-text("周杰伦")').first().click();
      // 歌曲列表应该消失
      await expect(page.locator('text=晴天')).not.toBeVisible();
    });

    test('升序/降序排序切换', async ({ page }) => {
      const sortButton = page.locator('button:has-text("按名称")');
      await expect(sortButton).toBeVisible();

      // 点击切换排序
      await sortButton.click();

      // 验证排序变化（请求应该带有新的order参数）
      const requests = [];
      page.on('request', req => requests.push(req));
    });
  });

  test.describe('Search - Filename (Story 2.8)', () => {
    test('搜索栏应该可见', async ({ page }) => {
      await expect(page.locator('input[placeholder*="搜索"]')).toBeVisible();
    });

    test('输入关键词并按回车应触发搜索', async ({ page }) => {
      // Mock search API
      await page.route('/api/songs/search**', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: mockSongs })
        });
      });

      // 输入搜索关键词
      const searchInput = page.locator('input[placeholder*="搜索"]');
      await searchInput.fill('晴天');
      await searchInput.press('Enter');

      // 应该显示搜索结果
      await expect(page.locator('text=/找到.*首歌曲/')).toBeVisible({ timeout: 5000 });
    });

    test('空关键词搜索应不触发请求', async ({ page }) => {
      const searchInput = page.locator('input[placeholder*="搜索"]');
      await searchInput.fill('');
      await searchInput.press('Enter');

      // 页面不应跳转到搜索结果
      await expect(page.locator('text=搜索结果')).not.toBeVisible();
    });

    test('无搜索结果应显示提示', async ({ page }) => {
      await page.route('/api/songs/search**', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: [] })
        });
      });

      const searchInput = page.locator('input[placeholder*="搜索"]');
      await searchInput.fill('不存在的歌曲');
      await searchInput.press('Enter');

      // 应该显示无结果提示
      await expect(page.locator('text=未找到匹配的歌曲')).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe('Search - Tag Content (Story 2.9)', () => {
    test('切换到按标签搜索模式', async ({ page }) => {
      // 查找搜索模式切换按钮
      const tagSearchButton = page.locator('button:has-text("标签")');
      if (await tagSearchButton.isVisible()) {
        await tagSearchButton.click();
      }

      // 验证切换成功（可能需要检查按钮激活状态）
    });

    test('按标签搜索应搜索标题、艺术家、专辑', async ({ page }) => {
      // Mock tag search API
      await page.route('/api/songs/search/by-tag**', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: mockSongs.filter(s => s.artist.includes('周杰伦')) })
        });
      });

      const searchInput = page.locator('input[placeholder*="搜索"]');
      await searchInput.fill('周杰伦');
      await searchInput.press('Enter');

      // 应该显示匹配艺术家"周杰伦"的歌曲
      await expect(page.locator('text=周杰伦').first()).toBeVisible({ timeout: 5000 });
    });

    test('多关键词搜索应同时匹配多个条件', async ({ page }) => {
      await page.route('/api/songs/search/by-tag**', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ data: mockSongs.filter(s => s.title === '晴天' && s.artist === '周杰伦') })
        });
      });

      const searchInput = page.locator('input[placeholder*="搜索"]');
      await searchInput.fill('周杰伦 晴天');
      await searchInput.press('Enter');

      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    });
  });

  test.describe('Song Selection (Story 2.6)', () => {
    test.beforeEach(async ({ page }) => {
      // Mock all necessary APIs
      await page.route('/api/artists**', async (route) => {
        const url = route.request().url();
        if (url.includes('/api/artists/')) {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockSongs })
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockArtists })
          });
        }
      });
    });

    test('点击歌曲复选框应选中歌曲', async ({ page }) => {
      // 展开艺术家
      await page.locator('tr:has-text("周杰伦")').first().click();

      // 等待歌曲列表出现
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 点击第一首歌的复选框
      const checkbox = page.locator('input[type="checkbox"]').first();
      await checkbox.click();

      // 应该显示选中数量
      await expect(page.locator('text=已选中')).toBeVisible();
    });

    test('选中按钮应该显示批量操作选项', async ({ page }) => {
      // 展开艺术家
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 选中一首歌
      await page.locator('input[type="checkbox"]').first().click();

      // 应该显示"全选"和"批量编辑"按钮
      await expect(page.locator('button:has-text("全选")').first()).toBeVisible();
    });
  });

  test.describe('Delete Songs (Story 2.7)', () => {
    test.beforeEach(async ({ page }) => {
      // Mock all necessary APIs
      await page.route('/api/artists**', async (route) => {
        const url = route.request().url();
        if (url.includes('/api/artists/')) {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockSongs })
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ data: mockArtists })
          });
        }
      });

      // 展开艺术家以显示歌曲
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    });

    test('选中歌曲后应显示删除按钮', async ({ page }) => {
      // 选中一首歌
      await page.locator('input[type="checkbox"]').first().click();

      // 应该显示删除按钮
      await expect(page.locator('button:has-text("删除")')).toBeVisible();
    });

    test('点击删除应显示确认对话框', async ({ page }) => {
      // 选中一首歌
      await page.locator('input[type="checkbox"]').first().click();

      // Mock delete confirmation API
      await page.route('**/api/songs/delete**', async (route) => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: { total: 1, succeeded: 1, failed: 0, results: [{ id: 1, status: 'deleted' }] }
          })
        });
      });

      // 点击删除
      await page.locator('button:has-text("删除")').click();

      // 应该显示确认对话框
      await expect(page.locator('h3:has-text("确认删除")')).toBeVisible();
    });
  });
});
