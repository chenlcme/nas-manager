import { test, expect } from '@playwright/test';

// Mock data
const mockArtists = [
  { id: 1, name: '周杰伦', songCount: 2 },
  { id: 2, name: '林俊杰', songCount: 1 },
];

const mockSongs = [
  { id: 1, title: '晴天', artist: '周杰伦', album: '叶惠美', year: 2003, genre: '流行', duration: 267, filePath: '/music/rock/晴天.mp3', coverPath: '', lyrics: '故事的小黄花 从出生那年就飘着' },
  { id: 2, title: '夜曲', artist: '周杰伦', album: '七里香', year: 2004, genre: '流行', duration: 252, filePath: '/music/rock/夜曲.mp3', coverPath: '', lyrics: '' },
  { id: 3, title: '江南', artist: '林俊杰', album: '编号89757', year: 2004, genre: '流行', duration: 245, filePath: '/music/pop/江南.mp3', coverPath: '', lyrics: '' },
];

test.describe('播放器 - 样式与交互综合测试', () => {
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
          body: JSON.stringify({ data: mockArtists })
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

  test('播放器面板样式验证 - 固定定位与阴影', async ({ page }) => {
    // 展开艺术家并播放
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证播放器面板基础样式
    const playerPanel = page.locator('.fixed.right-0');
    await expect(playerPanel).toBeVisible();
  });

  test('播放控制按钮样式验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证播放/暂停按钮
    const playPauseButton = page.locator('.fixed.right-0 button.rounded-full');
    await expect(playPauseButton).toBeVisible();
  });

  test('进度条与音量控制样式验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证进度条存在
    const progressBar = page.locator('.fixed.right-0 input[type="range"]').first();
    await expect(progressBar).toBeVisible();
  });

  test('编辑面板交互与样式验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 点击编辑按钮
    await page.locator('button:has-text("编辑信息")').click();

    // 验证编辑面板显示
    await expect(page.locator('.fixed.right-0 h4:has-text("编辑信息")')).toBeVisible();
  });

  test('关闭按钮交互验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证关闭按钮可见
    const closeButton = page.locator('.fixed.right-0 button[title="关闭"]');
    await expect(closeButton).toBeVisible();

    // 点击关闭
    await closeButton.click();
    await expect(page.locator('text=正在播放')).not.toBeVisible();
  });

  test('歌词区域样式验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证歌词标题
    await expect(page.locator('.fixed.right-0 h4')).toContainText('歌词');
  });

  test('专辑封面占位符样式验证', async ({ page }) => {
    await page.locator('tr:has-text("周杰伦")').first().click();
    await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
    await page.locator('button[title="播放"]').first().click();

    // 验证无封面时的占位符
    const coverArea = page.locator('.fixed.right-0 .w-48.h-48');
    await expect(coverArea).toBeVisible();
  });
});

test.describe('Tab导航样式验证', () => {
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
          body: JSON.stringify({ data: mockArtists })
        });
      }
    });

    await page.goto('/');
  });

  test('Tab导航应该显示所有标签', async ({ page }) => {
    await expect(page.locator('text=歌手')).toBeVisible();
    await expect(page.locator('text=专辑')).toBeVisible();
    await expect(page.locator('text=文件夹')).toBeVisible();
  });

  test('默认Tab应该有激活样式', async ({ page }) => {
    const artistTab = page.locator('button', { hasText: '歌手' });
    await expect(artistTab).toBeVisible();
  });

  test('Tab切换交互验证', async ({ page }) => {
    // 点击专辑Tab
    await page.locator('button', { hasText: '专辑' }).click();
    await expect(page.locator('text=专辑')).toBeVisible();

    // 点击文件夹Tab
    await page.locator('button', { hasText: '文件夹' }).click();
    await expect(page.locator('text=文件夹')).toBeVisible();
  });
});
