# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: epic2-browse-search.spec.ts >> Epic 2: 音乐库浏览与搜索 >> Delete Songs (Story 2.7) >> 点击删除应显示确认对话框
- Location: e2e/epic2-browse-search.spec.ts:321:5

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('text=确认删除')
Expected: visible
Error: strict mode violation: locator('text=确认删除') resolved to 2 elements:
    1) <h3 class="text-lg font-semibold text-gray-900">确认删除</h3> aka getByRole('heading', { name: '确认删除' })
    2) <button class="px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500">确认删除</button> aka getByRole('button', { name: '确认删除' })

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('text=确认删除')

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
          - row "已选中 1 / 3 首歌曲 全选 取消全选 批量编辑 删除 确认删除 将删除 1 首歌曲，文件和数据库记录都将被删除。 此操作不可撤销 取消 确认删除 按名称 名称 歌名 艺术家 专辑 年份 流派 时长 无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情 无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情 无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e34]:
            - cell "已选中 1 / 3 首歌曲 全选 取消全选 批量编辑 删除 确认删除 将删除 1 首歌曲，文件和数据库记录都将被删除。 此操作不可撤销 取消 确认删除 按名称 名称 歌名 艺术家 专辑 年份 流派 时长 无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情 无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情 无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e35]:
              - generic [ref=e37]:
                - generic [ref=e38]:
                  - generic [ref=e39]: 已选中 1 / 3 首歌曲
                  - button "全选" [ref=e40]
                  - button "取消全选" [ref=e41]
                - generic [ref=e42]:
                  - button "批量编辑" [ref=e43]
                  - button "删除" [active] [ref=e44]
              - generic [ref=e46]:
                - generic [ref=e47]:
                  - generic [ref=e48]:
                    - img [ref=e50]
                    - heading "确认删除" [level=3] [ref=e52]
                  - generic [ref=e53]:
                    - paragraph [ref=e54]: 将删除 1 首歌曲，文件和数据库记录都将被删除。
                    - paragraph [ref=e55]: 此操作不可撤销
                - generic [ref=e56]:
                  - button "取消" [ref=e57]
                  - button "确认删除" [ref=e58]
              - generic [ref=e60]:
                - combobox [ref=e61]:
                  - option "按名称" [selected]
                  - option "按时长"
                  - option "按添加时间"
                - button "名称" [ref=e62]:
                  - text: 名称
                  - img [ref=e64]
              - table [ref=e66]:
                - rowgroup [ref=e67]:
                  - row "歌名 艺术家 专辑 年份 流派 时长" [ref=e68]:
                    - columnheader [ref=e69]
                    - columnheader [ref=e70]
                    - columnheader "歌名" [ref=e71]
                    - columnheader "艺术家" [ref=e72]
                    - columnheader "专辑" [ref=e73]
                    - columnheader "年份" [ref=e74]
                    - columnheader "流派" [ref=e75]
                    - columnheader "时长" [ref=e76]
                    - columnheader [ref=e77]
                - rowgroup [ref=e78]:
                  - row "无封面 晴天 周杰伦 叶惠美 2003 流行 4:27 播放歌曲 查看歌曲详情" [ref=e79]:
                    - cell [ref=e80]:
                      - checkbox [checked] [ref=e81]
                    - cell "无封面" [ref=e82]:
                      - generic [ref=e83]: 无封面
                    - cell "晴天" [ref=e84]:
                      - generic [ref=e85]: 晴天
                    - cell "周杰伦" [ref=e86]
                    - cell "叶惠美" [ref=e87]
                    - cell "2003" [ref=e88]
                    - cell "流行" [ref=e89]
                    - cell "4:27" [ref=e90]
                    - cell "播放歌曲 查看歌曲详情" [ref=e91]:
                      - generic [ref=e92]:
                        - button "播放歌曲" [ref=e93]:
                          - img [ref=e94]
                        - button "查看歌曲详情" [ref=e96]:
                          - img [ref=e97]
                  - row "无封面 夜曲 周杰伦 七里香 2004 流行 4:12 播放歌曲 查看歌曲详情" [ref=e99]:
                    - cell [ref=e100]:
                      - checkbox [ref=e101]
                    - cell "无封面" [ref=e102]:
                      - generic [ref=e103]: 无封面
                    - cell "夜曲" [ref=e104]:
                      - generic [ref=e105]: 夜曲
                    - cell "周杰伦" [ref=e106]
                    - cell "七里香" [ref=e107]
                    - cell "2004" [ref=e108]
                    - cell "流行" [ref=e109]
                    - cell "4:12" [ref=e110]
                    - cell "播放歌曲 查看歌曲详情" [ref=e111]:
                      - generic [ref=e112]:
                        - button "播放歌曲" [ref=e113]:
                          - img [ref=e114]
                        - button "查看歌曲详情" [ref=e116]:
                          - img [ref=e117]
                  - row "无封面 江南 林俊杰 编号89757 2004 流行 4:05 播放歌曲 查看歌曲详情" [ref=e119]:
                    - cell [ref=e120]:
                      - checkbox [ref=e121]
                    - cell "无封面" [ref=e122]:
                      - generic [ref=e123]: 无封面
                    - cell "江南" [ref=e124]:
                      - generic [ref=e125]: 江南
                    - cell "林俊杰" [ref=e126]
                    - cell "编号89757" [ref=e127]
                    - cell "2004" [ref=e128]
                    - cell "流行" [ref=e129]
                    - cell "4:05" [ref=e130]
                    - cell "播放歌曲 查看歌曲详情" [ref=e131]:
                      - generic [ref=e132]:
                        - button "播放歌曲" [ref=e133]:
                          - img [ref=e134]
                        - button "查看歌曲详情" [ref=e136]:
                          - img [ref=e137]
          - row "› 林俊杰 2 首歌曲" [ref=e139]:
            - cell "›" [ref=e140]
            - cell "林俊杰" [ref=e141]
            - cell "2 首歌曲" [ref=e142]
```

# Test source

```ts
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
  272 |       await expect(page.locator('text=已选择')).toBeVisible();
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
> 340 |       await expect(page.locator('text=确认删除')).toBeVisible();
      |                                               ^ Error: expect(locator).toBeVisible() failed
  341 |     });
  342 |   });
  343 | });
  344 | 
```