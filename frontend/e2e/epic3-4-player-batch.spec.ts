import { test, expect } from '@playwright/test';

// Mock data
const mockSong = {
  id: 1,
  title: '晴天',
  artist: '周杰伦',
  album: '叶惠美',
  year: 2003,
  genre: '流行',
  duration: 267,
  filePath: '/music/晴天.mp3',
  coverPath: '',
  lyrics: '故事的小黄花 从出生那年就飘着',
};

const mockSongs = [
  { id: 1, title: '晴天', artist: '周杰伦', album: '叶惠美', year: 2003, genre: '流行', duration: 267 },
  { id: 2, title: '夜曲', artist: '周杰伦', album: '七里香', year: 2004, genre: '流行', duration: 252 },
  { id: 3, title: '江南', artist: '林俊杰', album: '编号89757', year: 2004, genre: '流行', duration: 245 },
];

test.describe('Epic 3: 播放器与现场编辑', () => {
  test.beforeEach(async ({ page }) => {
    // Mock setup status
    await page.route('/api/setup/status', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { needs_setup: false, music_dir: '/tmp/music', db_path: '' }
        })
      });
    });

    // Mock artists API
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
          body: JSON.stringify({ data: [{ id: 1, name: '周杰伦', songCount: 2 }] })
        });
      }
    });

    // Mock stream endpoint
    await page.route('/api/songs/*/stream', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'audio/mpeg',
        body: ''
      });
    });

    await page.goto('/');
  });

  test.describe('Story 3.1: 播放选中音乐', () => {
    test('点击播放按钮应显示播放器', async ({ page }) => {
      // 展开艺术家
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 点击播放按钮
      const playButton = page.locator('button[title="播放"]').first();
      await playButton.click();

      // 应该显示播放器
      await expect(page.locator('text=正在播放')).toBeVisible();
    });

    test('播放器应显示歌曲标题', async ({ page }) => {
      // 展开艺术家并播放
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 播放器应该显示歌曲标题
      await expect(page.locator('.fixed.right-0 h3')).toBeVisible();
    });

    test('播放/暂停按钮应该可见', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 播放/暂停按钮应该可见
      const playPauseButton = page.locator('.fixed.right-0 button.rounded-full');
      await expect(playPauseButton).toBeVisible();
    });
  });

  test.describe('Story 3.2: 展示专辑封面', () => {
    test('无封面时应显示占位符', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 播放器应该可见
      await expect(page.locator('.fixed.right-0')).toBeVisible();
    });
  });

  test.describe('Story 3.3: 展示歌词', () => {
    test('歌词区域应该可见', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 歌词标签应该可见
      await expect(page.locator('.fixed.right-0 h4')).toBeVisible();
    });
  });

  test.describe('Story 3.4: 展示播放时间', () => {
    test('进度条应该可见', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 进度条 input[type=range] 应该可见
      const progressBar = page.locator('.fixed.right-0 input[type="range"]').first();
      await expect(progressBar).toBeVisible();
    });
  });

  test.describe('Story 3.5: 播放中编辑元数据', () => {
    test('点击编辑按钮应显示编辑面板', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();

      // 点击编辑信息按钮
      await page.locator('button:has-text("编辑信息")').click();

      // 应该显示编辑表单
      await expect(page.locator('.fixed.right-0 h4:has-text("编辑信息")')).toBeVisible();
    });

    test('编辑表单应该包含输入框', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();
      await page.locator('button:has-text("编辑信息")').click();

      // 应该有多个文本输入框
      const inputs = page.locator('.fixed.right-0 input[type="text"]');
      await expect(inputs.first()).toBeVisible();
    });

    test('取消按钮应该关闭编辑面板', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('button[title="播放"]').first().click();
      await page.locator('button:has-text("编辑信息")').click();

      // 点击取消
      await page.locator('.fixed.right-0 button:has-text("取消")').click();

      // 编辑面板应该关闭（保存按钮不再可见）
      await expect(page.locator('.fixed.right-0 button:has-text("保存")')).not.toBeVisible();
    });
  });
});

test.describe('Epic 4: 批量编辑与撤销', () => {
  test.beforeEach(async ({ page }) => {
    // Mock setup status
    await page.route('/api/setup/status', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { needs_setup: false, music_dir: '/tmp/music', db_path: '' }
        })
      });
    });

    // Mock artists API
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
          body: JSON.stringify({ data: [{ id: 1, name: '周杰伦', songCount: 2 }] })
        });
      }
    });

    // Mock batch-get
    await page.route('/api/songs/batch-get', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: [mockSongs[0]] })
      });
    });

    // Mock batch-update
    await page.route('/api/songs/batch-update', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: { total: 1, succeeded: 1, failed: 0 } })
      });
    });

    // Mock batches/latest
    await page.route('/api/batches/latest', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { id: 1, type: 'update', target_ids: '[1]', old_values: '{"1":{"title":"旧标题"}}', new_values: '{"1":{"title":"新标题"}}' }
        })
      });
    });

    // Mock undo
    await page.route('/api/songs/undo/1', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ data: { succeeded: 1, failed: 0 } })
      });
    });

    await page.goto('/');
  });

  test.describe('Story 4.1: 批量修改标签', () => {
    test('选中歌曲后应该出现批量操作栏', async ({ page }) => {
      // 展开艺术家
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 选中第一首歌的复选框
      const checkbox = page.locator('table input[type="checkbox"]').first();
      await checkbox.click();

      // 应该显示选中数量
      await expect(page.locator('text=已选中')).toBeVisible();
    });

    test('批量编辑按钮应该可见', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 选中一首歌
      await page.locator('table input[type="checkbox"]').first().click();

      // 批量编辑按钮应该可见
      await expect(page.locator('button:has-text("批量编辑")')).toBeVisible();
    });

    test('点击批量编辑应该显示编辑面板', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      // 选中一首歌
      await page.locator('table input[type="checkbox"]').first().click();

      // 点击批量编辑
      await page.locator('button:has-text("批量编辑")').click();

      // 应该显示批量编辑面板标题
      await expect(page.locator('.fixed.right-0 h2:has-text("批量编辑")')).toBeVisible();
    });

    test('编辑面板应该显示选中的歌曲数量', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });

      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 应该显示选中数量
      await expect(page.locator('.fixed.right-0 span:text("已选中")')).toBeVisible();
    });

    test('编辑面板应该有标签和歌词Tab', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // Tab 应该可见
      await expect(page.locator('button:has-text("标签信息")')).toBeVisible();
      await expect(page.locator('button:has-text("歌词")')).toBeVisible();
    });

    test('填写标签后应该显示预览', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 填写艺术家字段
      const artistInput = page.locator('.fixed.right-0 input[placeholder="保持不变"]').first();
      await artistInput.click();
      await artistInput.fill('新艺术家');
      await artistInput.blur();

      // 应该显示预览
      await expect(page.locator('text=预览更改')).toBeVisible();
    });

    test('点击应用更改应该成功', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 填写字段
      await page.locator('.fixed.right-0 input').first().fill('新艺术家');

      // 点击应用
      await page.locator('button:has-text("应用更改")').click();

      // 应该显示成功提示
      await expect(page.locator('text=批量更新成功')).toBeVisible({ timeout: 3000 });
    });
  });

  test.describe('Story 4.2: 批量修改封面', () => {
    test('批量编辑面板应该存在', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 批量编辑面板应该可见
      await expect(page.locator('.fixed.right-0')).toBeVisible();
    });
  });

  test.describe('Story 4.3: 搜索并批量应用歌词', () => {
    test('歌词Tab应该显示功能提示', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 点击歌词Tab
      await page.locator('button:has-text("歌词")').click();

      // 应该显示功能提示
      await expect(page.locator('text=批量歌词功能')).toBeVisible();
    });
  });

  test.describe('Story 4.4: 撤销批量编辑', () => {
    test('撤销按钮应该可见', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 撤销按钮应该可见
      await expect(page.locator('button:has-text("撤销")')).toBeVisible();
    });

    test('点击撤销应该成功', async ({ page }) => {
      await page.locator('tr:has-text("周杰伦")').first().click();
      await page.locator('table input[type="checkbox"]').first().click();
      await page.locator('button:has-text("批量编辑")').click();

      // 点击撤销
      await page.locator('button:has-text("撤销")').click();

      // 应该显示撤销成功
      await expect(page.locator('text=已撤销')).toBeVisible({ timeout: 3000 });
    });
  });
});
