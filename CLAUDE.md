# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git 操作规则

1. **禁止自动提交和推送**
   - 除非用户明确要求，否则不要自动执行 `git commit` 或 `git push`
   - 所有提交和推送操作必须经过用户确认

2. **Commit Message 按照语义化提交规范**
   - 使用以下格式：`type(scope): subject`
   - Type 类型：
     - `feat`: 新功能
     - `fix`: 修复 bug
     - `docs`: 文档变更
     - `style`: 代码格式调整（不影响功能）
     - `refactor`: 重构（不影响功能）
     - `test`: 测试相关
     - `chore`: 构建/工具变更
   - Scope：影响的模块或功能（如 `repo`, `handler`, `frontend` 等）
   - Subject：简短描述（不超过 50 字符，祈使语气）

   示例：
   ```
   feat(handler): 添加专辑视图浏览功能

   实现 Story 2.2 的后端 Handler，包括：
   - GET /api/albums 端点
   - GET /api/albums/:id/songs 端点
   ```
