# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: epic2-browse-search.spec.ts >> Epic 2: 音乐库浏览与搜索 >> Song Selection (Story 2.6) >> 点击歌曲复选框应选中歌曲
- Location: e2e/epic2-browse-search.spec.ts:260:5

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('text=已选择')
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('text=已选择')

```

# Page snapshot

```yaml
- generic [ref=e3]:
  - generic [ref=e4]:
    - generic [ref=e5]:
      - button "歌手" [ref=e6]
      - button "专辑" [ref=e7]
      - button "文件夹" [ref=e8]
    - generic [ref=e10]:
      - combobox [ref=e11]:
        - option "标签" [selected]
        - option "文件名"
      - generic [ref=e12]:
        - textbox "搜索标题、艺术家、专辑..." [ref=e13]
        - img [ref=e14]
      - button "搜索" [disabled] [ref=e16]
  - main [ref=e17]:
    - generic [ref=e18]:
      - generic [ref=e19]:
        - generic [ref=e20]: 共 2 位艺术家
        - button "按名称降序 ↕" [ref=e21]
      - table [ref=e23]:
        - rowgroup [ref=e24]:
          - row "艺术家 歌曲数量" [ref=e25]:
            - columnheader [ref=e26]
            - columnheader "艺术家" [ref=e27]
            - columnheader "歌曲数量" [ref=e28]
        - rowgroup [ref=e29]:
          - row "› 周杰伦 3 首歌曲" [ref=e30]:
            - cell "›" [ref=e31]
            - cell "周杰伦" [ref=e32]
            - cell "3 首歌曲" [ref=e33]
          - row "已选中 1 / 3 首歌曲 全选 取消全选 批量编辑 删除 按名称 名称 歌名 艺术家 专辑 年份 流派 时长 无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情 无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情 无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e34]:
            - cell "已选中 1 / 3 首歌曲 全选 取消全选 批量编辑 删除 按名称 名称 歌名 艺术家 专辑 年份 流派 时长 无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情 无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情 无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e35]:
              - generic [ref=e37]:
                - generic [ref=e38]:
                  - generic [ref=e39]: 已选中 1 / 3 首歌曲
                  - button "全选" [ref=e40]
                  - button "取消全选" [ref=e41]
                - generic [ref=e42]:
                  - button "批量编辑" [ref=e43]
                  - button "删除" [ref=e44]
              - generic [ref=e46]:
                - combobox [ref=e47]:
                  - option "按名称" [selected]
                  - option "按时长"
                  - option "按添加时间"
                - button "名称" [ref=e48]:
                  - text: 名称
                  - img [ref=e50]
              - table [ref=e52]:
                - rowgroup [ref=e53]:
                  - row "歌名 艺术家 专辑 年份 流派 时长" [ref=e54]:
                    - columnheader [ref=e55]
                    - columnheader [ref=e56]
                    - columnheader "歌名" [ref=e57]
                    - columnheader "艺术家" [ref=e58]
                    - columnheader "专辑" [ref=e59]
                    - columnheader "年份" [ref=e60]
                    - columnheader "流派" [ref=e61]
                    - columnheader "时长" [ref=e62]
                    - columnheader [ref=e63]
                - rowgroup [ref=e64]:
                  - row "无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情" [ref=e65]:
                    - cell [ref=e66]:
                      - checkbox [checked] [active] [ref=e67]
                    - cell "无封面" [ref=e68]:
                      - generic [ref=e69]: 无封面
                    - cell "晴天" [ref=e70]:
                      - generic [ref=e71]: 晴天
                    - cell "周杰伦" [ref=e72]
                    - cell "叶惠美" [ref=e73]
                    - cell "2003" [ref=e74]
                    - cell "流行" [ref=e75]
                    - cell "4:27" [ref=e76]
                    - cell "播放歌曲 查看歌曲详情" [ref=e77]:
                      - generic [ref=e78]:
                        - button "播放歌曲" [ref=e79]:
                          - img [ref=e80]
                        - button "查看歌曲详情" [ref=e82]:
                          - img [ref=e83]
                  - row "无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情" [ref=e85]:
                    - cell [ref=e86]:
                      - checkbox [ref=e87]
                    - cell "无封面" [ref=e88]:
                      - generic [ref=e89]: 无封面
                    - cell "夜曲" [ref=e90]:
                      - generic [ref=e91]: 夜曲
                    - cell "周杰伦" [ref=e92]
                    - cell "七里香" [ref=e93]
                    - cell "2004" [ref=e94]
                    - cell "流行" [ref=e95]
                    - cell "4:12" [ref=e96]
                    - cell "播放歌曲 查看歌曲详情" [ref=e97]:
                      - generic [ref=e98]:
                        - button "播放歌曲" [ref=e99]:
                          - img [ref=e100]
                        - button "查看歌曲详情" [ref=e102]:
                          - img [ref=e103]
                  - row "无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e105]:
                    - cell [ref=e106]:
                      - checkbox [ref=e107]
                    - cell "无封面" [ref=e108]:
                      - generic [ref=e109]: 无封面
                    - cell "江南" [ref=e110]:
                      - generic [ref=e111]: 江南
                    - cell "林俊杰" [ref=e112]
                    - cell "编号89757" [ref=e113]
                    - cell "2004" [ref=e114]
                    - cell "流行" [ref=e115]
                    - cell "4:05" [ref=e116]
                    - cell "播放歌曲 查看歌曲详情" [ref=e117]:
                      - generic [ref=e118]:
                        - button "播放歌曲" [ref=e119]:
                          - img [ref=e120]
                        - button "查看歌曲详情" [ref=e122]:
                          - img [ref=e123]
          - row "› 林俊杰 2 首歌曲" [ref=e125]:
            - cell "›" [ref=e126]
            - cell "林俊杰" [ref=e127]
            - cell "2 首歌曲" [ref=e128]
```

# Test source

```ts
  172 |       await expect(page.locator('text=搜索结果')).not.toBeVisible();
  173 |     });
  174 | 
  175 |     test('无搜索结果应显示提示', async ({ page }) => {
  176 |       await page.route('/api/songs/search**', async (route) => {
  177 |         await route.fulfill({
  178 |           status: 200,
  179 |           contentType: 'application/json',
  180 |           body: JSON.stringify({ data: [] })
  181 |         });
  182 |       });
  183 | 
  184 |       const searchInput = page.locator('input[placeholder*="搜索"]');
  185 |       await searchInput.fill('不存在的歌曲');
  186 |       await searchInput.press('Enter');
  187 | 
  188 |       // 应该显示无结果提示
  189 |       await expect(page.locator('text=未找到匹配的歌曲')).toBeVisible({ timeout: 5000 });
  190 |     });
  191 |   });
  192 | 
  193 |   test.describe('Search - Tag Content (Story 2.9)', () => {
  194 |     test('切换到按标签搜索模式', async ({ page }) => {
  195 |       // 查找搜索模式切换按钮
  196 |       const tagSearchButton = page.locator('button:has-text("标签")');
  197 |       if (await tagSearchButton.isVisible()) {
  198 |         await tagSearchButton.click();
  199 |       }
  200 | 
  201 |       // 验证切换成功（可能需要检查按钮激活状态）
  202 |     });
  203 | 
  204 |     test('按标签搜索应搜索标题、艺术家、专辑', async ({ page }) => {
  205 |       // Mock tag search API
  206 |       await page.route('/api/songs/search/by-tag**', async (route) => {
  207 |         await route.fulfill({
  208 |           status: 200,
  209 |           contentType: 'application/json',
  210 |           body: JSON.stringify({ data: mockSongs.filter(s => s.artist.includes('周杰伦')) })
  211 |         });
  212 |       });
  213 | 
  214 |       const searchInput = page.locator('input[placeholder*="搜索"]');
  215 |       await searchInput.fill('周杰伦');
  216 |       await searchInput.press('Enter');
  217 | 
  218 |       // 应该显示匹配艺术家"周杰伦"的歌曲
  219 |       await expect(page.locator('text=周杰伦').first()).toBeVisible({ timeout: 5000 });
  220 |     });
  221 | 
  222 |     test('多关键词搜索应同时匹配多个条件', async ({ page }) => {
  223 |       await page.route('/api/songs/search/by-tag**', async (route) => {
  224 |         await route.fulfill({
  225 |           status: 200,
  226 |           contentType: 'application/json',
  227 |           body: JSON.stringify({ data: mockSongs.filter(s => s.title === '晴天' && s.artist === '周杰伦') })
  228 |         });
  229 |       });
  230 | 
  231 |       const searchInput = page.locator('input[placeholder*="搜索"]');
  232 |       await searchInput.fill('周杰伦 晴天');
  233 |       await searchInput.press('Enter');
  234 | 
  235 |       await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
  236 |     });
  237 |   });
  238 | 
  239 |   test.describe('Song Selection (Story 2.6)', () => {
  240 |     test.beforeEach(async ({ page }) => {
  241 |       // Mock all necessary APIs
  242 |       await page.route('/api/artists**', async (route) => {
  243 |         const url = route.request().url();
  244 |         if (url.includes('/api/artists/')) {
  245 |           await route.fulfill({
  246 |             status: 200,
  247 |             contentType: 'application/json',
  248 |             body: JSON.stringify({ data: mockSongs })
  249 |           });
  250 |         } else {
  251 |           await route.fulfill({
  252 |             status: 200,
  253 |             contentType: 'application/json',
  254 |             body: JSON.stringify({ data: mockArtists })
  255 |           });
  256 |         }
  257 |       });
  258 |     });
  259 | 
  260 |     test('点击歌曲复选框应选中歌曲', async ({ page }) => {
  261 |       // 展开艺术家
  262 |       await page.locator('tr:has-text("周杰伦")').first().click();
  263 | 
  264 |       // 等待歌曲列表出现
  265 |       await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
  266 | 
  267 |       // 点击第一首歌的复选框
  268 |       const checkbox = page.locator('input[type="checkbox"]').first();
  269 |       await checkbox.click();
  270 | 
  271 |       // 应该显示选中数量
> 272 |       await expect(page.locator('text=已选择')).toBeVisible();
      |                                              ^ Error: expect(locator).toBeVisible() failed
  273 |     });
  274 | 
  275 |     test('选中按钮应该显示批量操作选项', async ({ page }) => {
  276 |       // 展开艺术家
  277 |       await page.locator('tr:has-text("周杰伦")').first().click();
  278 |       await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
  279 | 
  280 |       // 选中一首歌
  281 |       await page.locator('input[type="checkbox"]').first().click();
  282 | 
  283 |       // 应该显示"全选"和"批量编辑"按钮
  284 |       await expect(page.locator('button:has-text("全选")')).toBeVisible();
  285 |     });
  286 |   });
  287 | 
  288 |   test.describe('Delete Songs (Story 2.7)', () => {
  289 |     test.beforeEach(async ({ page }) => {
  290 |       // Mock all necessary APIs
  291 |       await page.route('/api/artists**', async (route) => {
  292 |         const url = route.request().url();
  293 |         if (url.includes('/api/artists/')) {
  294 |           await route.fulfill({
  295 |             status: 200,
  296 |             contentType: 'application/json',
  297 |             body: JSON.stringify({ data: mockSongs })
  298 |           });
  299 |         } else {
  300 |           await route.fulfill({
  301 |             status: 200,
  302 |             contentType: 'application/json',
  303 |             body: JSON.stringify({ data: mockArtists })
  304 |           });
  305 |         }
  306 |       });
  307 | 
  308 |       // 展开艺术家以显示歌曲
  309 |       await page.locator('tr:has-text("周杰伦")').first().click();
  310 |       await expect(page.locator('text=晴天')).toBeVisible({ timeout: 5000 });
  311 |     });
  312 | 
  313 |     test('选中歌曲后应显示删除按钮', async ({ page }) => {
  314 |       // 选中一首歌
  315 |       await page.locator('input[type="checkbox"]').first().click();
  316 | 
  317 |       // 应该显示删除按钮
  318 |       await expect(page.locator('button:has-text("删除")')).toBeVisible();
  319 |     });
  320 | 
  321 |     test('点击删除应显示确认对话框', async ({ page }) => {
  322 |       // 选中一首歌
  323 |       await page.locator('input[type="checkbox"]').first().click();
  324 | 
  325 |       // Mock delete confirmation API
  326 |       await page.route('**/api/songs/delete**', async (route) => {
  327 |         await route.fulfill({
  328 |           status: 200,
  329 |           contentType: 'application/json',
  330 |           body: JSON.stringify({
  331 |             data: { total: 1, succeeded: 1, failed: 0, results: [{ id: 1, status: 'deleted' }] }
  332 |           })
  333 |         });
  334 |       });
  335 | 
  336 |       // 点击删除
  337 |       await page.locator('button:has-text("删除")').click();
  338 | 
  339 |       // 应该显示确认对话框
  340 |       await expect(page.locator('text=确认删除')).toBeVisible();
  341 |     });
  342 |   });
  343 | });
  344 | 
```