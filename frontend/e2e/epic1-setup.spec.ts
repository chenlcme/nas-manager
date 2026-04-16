import { test, expect } from '@playwright/test';

test.describe('Epic 1: 首次配置向导 (Setup Wizard)', () => {
  test.beforeEach(async ({ page }) => {
    // Mock the setup status API - needs_setup=true means setup required
    await page.route('/api/setup/status', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { needs_setup: true, music_dir: '', db_path: '' }
        })
      });
    });

    await page.goto('/');
  });

  test('应该显示首次配置向导', async ({ page }) => {
    // 验证标题
    await expect(page.locator('h1')).toContainText('首次配置向导');

    // 验证步骤指示器显示步骤1
    await expect(page.locator('h2')).toContainText('步骤 1：设置音乐目录');

    // 验证有音乐目录输入框
    const musicDirInput = page.locator('#music-dir');
    await expect(musicDirInput).toBeVisible();

    // 验证有"下一步"按钮
    const nextButton = page.locator('button:has-text("下一步")');
    await expect(nextButton).toBeVisible();
  });

  test('空目录路径应显示错误提示', async ({ page }) => {
    // 点击下一步而不输入路径
    const nextButton = page.locator('button:has-text("下一步")');
    await nextButton.click();

    // 应该显示错误信息
    await expect(page.locator('.text-red-500')).toContainText('请输入音乐目录路径');
  });

  test('应能进入步骤2', async ({ page }) => {
    // 输入有效的目录路径
    const musicDirInput = page.locator('#music-dir');
    await musicDirInput.fill('/tmp/test-music');

    // 点击下一步
    const nextButton = page.locator('button:has-text("下一步")');
    await nextButton.click();

    // 应该显示步骤2
    await expect(page.locator('h2')).toContainText('步骤 2：确认数据库路径');

    // 验证有数据库路径输入框
    const dbPathInput = page.locator('#db-path');
    await expect(dbPathInput).toBeVisible();

    // 验证有"返回"和"完成配置"按钮
    await expect(page.locator('button:has-text("返回")')).toBeVisible();
    await expect(page.locator('button:has-text("完成配置")')).toBeVisible();
  });

  test('应能从步骤2返回步骤1', async ({ page }) => {
    // 先进入步骤2
    const musicDirInput = page.locator('#music-dir');
    await musicDirInput.fill('/tmp/test-music');
    await page.locator('button:has-text("下一步")').click();

    // 点击返回
    await page.locator('button:has-text("返回")').click();

    // 应该回到步骤1
    await expect(page.locator('h2')).toContainText('步骤 1：设置音乐目录');
  });

  test('不存在的目录应显示错误', async ({ page, baseURL }) => {
    // Mock API to return success after validation
    await page.route('/api/setup', async (route) => {
      await route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({
          error: {
            code: 'DIR_NOT_EXIST',
            message: 'Music directory does not exist'
          }
        })
      });
    });

    // 输入不存在的目录
    const musicDirInput = page.locator('#music-dir');
    await musicDirInput.fill('/non/existent/path');
    await page.locator('button:has-text("下一步")').click();

    // 进入步骤2
    await expect(page.locator('h2')).toContainText('步骤 2：确认数据库路径');

    // 点击完成配置
    await page.locator('button:has-text("完成配置")').click();

    // 应该显示错误
    await expect(page.locator('.text-red-500')).toBeVisible();
  });

  test('配置成功后应跳转到主页面', async ({ page }) => {
    // Mock API to return success
    await page.route('/api/setup', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: { success: true }
        })
      });
    });

    // 输入有效目录并进入步骤2
    const musicDirInput = page.locator('#music-dir');
    await musicDirInput.fill('/tmp/test-music');
    await page.locator('button:has-text("下一步")').click();

    // 点击完成配置
    await page.locator('button:has-text("完成配置")').click();

    // 等待页面跳转或内容变化
    await page.waitForURL('**/');
  });
});
