# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: epic2-browse-search.spec.ts >> Epic 2: 音乐库浏览与搜索 >> Search - Filename (Story 2.8) >> 输入关键词并按回车应触发搜索
- Location: e2e/epic2-browse-search.spec.ts:147:5

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('text=搜索结果')
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('text=搜索结果')

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
      - combobox [disabled] [ref=e11]:
        - option "标签" [selected]
        - option "文件名"
      - generic [ref=e12]:
        - textbox "搜索标题、艺术家、专辑..." [disabled] [ref=e13]: 晴天
        - img [ref=e14]
      - button "搜索" [disabled] [ref=e16]
  - main [ref=e17]:
    - generic [ref=e18]:
      - generic [ref=e19]:
        - button "返回" [ref=e20]:
          - img [ref=e21]
          - text: 返回
        - generic [ref=e23]: 搜索 "晴天" - 找到 3 首歌曲
      - table [ref=e25]:
        - rowgroup [ref=e26]:
          - row "歌曲 歌手 专辑 时长" [ref=e27]:
            - columnheader [ref=e28]
            - columnheader "歌曲" [ref=e29]
            - columnheader "歌手" [ref=e30]
            - columnheader "专辑" [ref=e31]
            - columnheader "时长" [ref=e32]
        - rowgroup [ref=e33]:
          - row "晴天 周杰伦 叶惠美 4:27" [ref=e34]:
            - cell [ref=e35]:
              - button [ref=e36]:
                - img [ref=e37]
            - cell "晴天" [ref=e39]:
              - mark [ref=e41]: 晴天
            - cell "周杰伦" [ref=e42]:
              - generic [ref=e43]: 周杰伦
            - cell "叶惠美" [ref=e44]:
              - generic [ref=e45]: 叶惠美
            - cell "4:27" [ref=e46]
          - row "夜曲 周杰伦 七里香 4:12" [ref=e47]:
            - cell [ref=e48]:
              - button [ref=e49]:
                - img [ref=e50]
            - cell "夜曲" [ref=e52]:
              - generic [ref=e53]: 夜曲
            - cell "周杰伦" [ref=e54]:
              - generic [ref=e55]: 周杰伦
            - cell "七里香" [ref=e56]:
              - generic [ref=e57]: 七里香
            - cell "4:12" [ref=e58]
          - row "江南 林俊杰 编号89757 4:05" [ref=e59]:
            - cell [ref=e60]:
              - button [ref=e61]:
                - img [ref=e62]
            - cell "江南" [ref=e64]:
              - generic [ref=e65]: 江南
            - cell "林俊杰" [ref=e66]:
              - generic [ref=e67]: 林俊杰
            - cell "编号89757" [ref=e68]:
              - generic [ref=e69]: 编号89757
            - cell "4:05" [ref=e70]
```

# Test source

```ts
  63  |       await page.route('/api/folders', async (route) => {
  64  |         await route.fulfill({
  65  |           status: 200,
  66  |           contentType: 'application/json',
  67  |           body: JSON.stringify({ data: [] })
  68  |         });
  69  |       });
  70  | 
  71  |       await page.locator('button:has-text("文件夹")').click();
  72  | 
  73  |       await expect(page.locator('text=文件夹')).toBeVisible();
  74  |     });
  75  |   });
  76  | 
  77  |   test.describe('Artist View (Story 2.1)', () => {
  78  |     test.beforeEach(async ({ page }) => {
  79  |       // Mock artists API
  80  |       await page.route('/api/artists**', async (route) => {
  81  |         const url = route.request().url();
  82  |         if (url.includes('/api/artists/')) {
  83  |           // Artist songs endpoint
  84  |           await route.fulfill({
  85  |             status: 200,
  86  |             contentType: 'application/json',
  87  |             body: JSON.stringify({ data: mockSongs.filter(s => s.artist === '周杰伦') })
  88  |           });
  89  |         } else {
  90  |           // Artists list endpoint
  91  |           await route.fulfill({
  92  |             status: 200,
  93  |             contentType: 'application/json',
  94  |             body: JSON.stringify({ data: mockArtists })
  95  |           });
  96  |         }
  97  |       });
  98  |     });
  99  | 
  100 |     test('应显示艺术家列表', async ({ page }) => {
  101 |       // 等待艺术家列表加载
  102 |       await expect(page.locator('text=周杰伦')).toBeVisible({ timeout: 5000 });
  103 |       await expect(page.locator('text=林俊杰')).toBeVisible();
  104 |     });
  105 | 
  106 |     test('应显示艺术家数量', async ({ page }) => {
  107 |       await expect(page.locator('text=共 2 位艺术家')).toBeVisible();
  108 |     });
  109 | 
  110 |     test('点击艺术家应展开歌曲列表', async ({ page }) => {
  111 |       // 点击周杰伦
  112 |       await page.locator('tr:has-text("周杰伦")').first().click();
  113 | 
  114 |       // 应该展开显示歌曲
  115 |       await expect(page.locator('text=晴天')).toBeVisible();
  116 |     });
  117 | 
  118 |     test('再次点击艺术家应折叠歌曲列表', async ({ page }) => {
  119 |       // 点击展开
  120 |       await page.locator('tr:has-text("周杰伦")').first().click();
  121 |       await expect(page.locator('text=晴天')).toBeVisible();
  122 | 
  123 |       // 点击折叠
  124 |       await page.locator('tr:has-text("周杰伦")').first().click();
  125 |       // 歌曲列表应该消失
  126 |       await expect(page.locator('text=晴天')).not.toBeVisible();
  127 |     });
  128 | 
  129 |     test('升序/降序排序切换', async ({ page }) => {
  130 |       const sortButton = page.locator('button:has-text("按名称")');
  131 |       await expect(sortButton).toBeVisible();
  132 | 
  133 |       // 点击切换排序
  134 |       await sortButton.click();
  135 | 
  136 |       // 验证排序变化（请求应该带有新的order参数）
  137 |       const requests = [];
  138 |       page.on('request', req => requests.push(req));
  139 |     });
  140 |   });
  141 | 
  142 |   test.describe('Search - Filename (Story 2.8)', () => {
  143 |     test('搜索栏应该可见', async ({ page }) => {
  144 |       await expect(page.locator('input[placeholder*="搜索"]')).toBeVisible();
  145 |     });
  146 | 
  147 |     test('输入关键词并按回车应触发搜索', async ({ page }) => {
  148 |       // Mock search API
  149 |       await page.route('/api/songs/search**', async (route) => {
  150 |         await route.fulfill({
  151 |           status: 200,
  152 |           contentType: 'application/json',
  153 |           body: JSON.stringify({ data: mockSongs })
  154 |         });
  155 |       });
  156 | 
  157 |       // 输入搜索关键词
  158 |       const searchInput = page.locator('input[placeholder*="搜索"]');
  159 |       await searchInput.fill('晴天');
  160 |       await searchInput.press('Enter');
  161 | 
  162 |       // 应该显示搜索结果
> 163 |       await expect(page.locator('text=搜索结果')).toBeVisible({ timeout: 5000 });
      |                                               ^ Error: expect(locator).toBeVisible() failed
  164 |     });
  165 | 
  166 |     test('空关键词搜索应不触发请求', async ({ page }) => {
  167 |       const searchInput = page.locator('input[placeholder*="搜索"]');
  168 |       await searchInput.fill('');
  169 |       await searchInput.press('Enter');
  170 | 
  171 |       // 页面不应跳转到搜索结果
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
```