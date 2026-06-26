---
name: phased-goal-roadmap
description: Create reusable phased /goal Markdown roadmaps with strict stage gates. Use when the user asks to write a goal command document, split work into execution phases, define phase acceptance criteria, enforce pre-commit multi-agent reviews, post-commit review loops, or browser screenshot validation for frontend changes.
---

# Phased Goal Roadmap

Use this skill to turn a user-defined requirement into a reusable `/goal` Markdown document. The user should only need to provide:

- the high-level requirement

The agent must infer sensible execution phases and phase acceptance criteria when the user does not provide them. The agent fills in the reusable structure, quality gates, review gates, browser validation gates, and final handoff rules.

## When To Use

Use this skill when the user says things like:

- "写一份 goal 命令可用的 md 文件"
- "结构参考某个 phased roadmap"
- "把需求拆成阶段和验收标准"
- "抽成可复用 skill"
- "提交前要多代理 review"
- "前端必须浏览器截图验收"

If the user has not provided enough requirement or phase detail, ask only the minimum critical questions. Do not over-specify implementation internals before reading the relevant code.

## Strict Format Gate

Any `/goal` prompt produced by this skill must use the canonical phased Markdown format below. If the user asks the agent to execute a `/goal` prompt that does not follow this format, the agent must refuse to execute it as a goal and first offer to convert it into the canonical format.

Reject or rewrite any goal prompt that is missing any of these required parts:

1. YAML frontmatter with `name`, `overview`, `todos`, and `isProject`.
2. A top-level title.
3. A `/goal 提示词` section explaining that the whole Markdown is the goal input.
4. `背景`.
5. `核心目标`.
6. `非目标`.
7. `总体原则`.
8. `强制阶段门禁`.
9. `Review 门禁细则`.
10. `浏览器截图验收门禁`.
11. One `Phase N` section per phase.
12. Each phase must contain exactly these subsections: `Codex Goal`, `Scope`, `Acceptance Criteria`.
13. `总体最终验收标准`.
14. `当前执行建议`.

The format must be visually recognizable as the same family as existing phased roadmap files such as `glossary-phased-roadmap.plan.md`: frontmatter first, then narrative context, then global gates, then numbered phases with explicit acceptance criteria.

If only a loose requirement is provided, do not reject it. Instead, create a canonical `/goal` Markdown file by inferring phases and acceptance criteria from the requirement and verified code context.

## Output Contract

Create or update a Markdown plan/goal document with:

1. YAML frontmatter containing `name`, `overview`, and phased `todos`.
2. A title and `/goal` usage note.
3. Background and current-state facts.
4. Core goals and non-goals.
5. Global principles.
6. Mandatory stage gates.
7. Review gates.
8. Browser screenshot validation gates.
9. Optional mermaid data flow if it clarifies architecture.
10. One section per phase, each with `Codex Goal`, `Scope`, and `Acceptance Criteria`.
11. Overall final acceptance criteria.
12. Current execution recommendation: only detail the next phase; later phases remain goal-level until the previous phase is complete.

Keep the document executable and clean. Do not include review debate history, model provenance, or implementation chatter.

## Phase Inference Rules

When the user provides only a requirement, infer phases using this order:

1. Foundation / data model / backend primitives.
2. Backend API or service integration.
3. Frontend entry, list, or shell if UI is involved.
4. Main user-facing workflow or detail view.
5. Notification, automation, migration, or operational integration.
6. End-to-end validation, documentation, and plan archival.

Collapse or expand phases based on scope. Prefer 3-6 phases. Each phase must be independently testable and committable.

For each phase, generate acceptance criteria that cover:

- the main happy path
- at least one relevant failure or empty state
- compatibility with existing behavior
- required automated tests
- browser screenshot validation when frontend changes are involved
- pre-commit five-subagent review
- post-commit review loop

## Per-Phase Planning Gate

Before executing any phase, the main agent must create an execution-level plan for that phase. The phase plan must be optimized before implementation starts.

Mandatory optimization workflow:

1. Read the canonical goal document and the current phase only.
2. Draft an execution plan for the current phase with concrete files, symbols, data flow, tests, browser validation needs, risks, and rollback boundaries.
3. Launch three adversarial plan-review subagents with self-contained review packets:
   - Builder reviewer: checks whether the plan is executable, complete, and sequenced correctly.
   - Skeptic reviewer: attacks assumptions, hidden coupling, over-design, under-design, and missing edge cases.
   - Verifier reviewer: verifies claims against code, docs, tests, and repository rules.
4. Run at least two optimization rounds. A round means: send the latest phase plan to all three reviewers, integrate evidence-backed feedback, then produce a revised phase plan.
5. If any reviewer finds a P0/P1 defect after round 2, continue additional rounds until no P0/P1 plan defect remains or report a blocker.
6. Only the optimized phase plan may enter implementation.

Do not include the review debate transcript in the final phase plan. Fold valid decisions into the clean execution plan and remove reviewer-history sections.

## Execution Delegation Model

The main agent is the workflow controller, not the default executor for every task.

During implementation:

- The main agent owns phase control, context packets, dependency ordering, integration decisions, conflict resolution, final validation, and user-facing status.
- Execution work should be delegated to subagents whenever it can be split by module, layer, or responsibility.
- Prefer multiple focused execution subagents over one broad "do everything" subagent.
- Each execution subagent must receive a self-contained packet: current phase goal, optimized phase plan, exact scope, allowed files or modules, constraints, tests to run, and expected output.
- Subagents should return concrete changes made, tests run, blockers, and risks. The main agent must inspect and integrate their outputs instead of trusting summaries blindly.
- The main agent may execute small glue tasks directly only when delegating would add more overhead than the task itself, or when a tool/action cannot be delegated safely.
- If subagents are unavailable or unsuitable, the main agent must state that constraint and continue only with explicit awareness that the normal delegation model is being bypassed.

## Required Global Gates

Every generated roadmap must include these gates unless the user explicitly removes them:

1. Create an execution-level phase plan.
2. Run at least two rounds of three-subagent adversarial plan optimization.
3. Delegate implementation work to focused execution subagents where practical, while the main agent controls integration and validation.
4. Complete the phase implementation.
5. Run phase-specific automated tests, type checks, syntax checks, and smoke checks.
6. If the phase touches frontend code, use browser tooling to open the real local page and take browser screenshots for visual validation.
7. Validate each phase acceptance criterion and record the result.
8. Before commit, run five specialized subagent reviews in parallel:
   - Security issues
   - Logic bugs
   - Test coverage
   - Maintainability
   - Performance / concurrency risks
9. If any specialized review finds an issue, fix it, rerun relevant tests and browser validation, then rerun all five specialized reviews. Do not rerun only the failed dimension.
10. Only when all five pre-commit reviews explicitly say no issues were found may the agent run commit-readiness checks.
11. Commit-readiness checks must include `git status`, `git diff`, project-required impact/change analysis such as GitNexus, and a check that credentials, generated artifacts, caches, and unrelated files are not included.
12. Create exactly one commit for the phase, using the repository's commit rules.
13. After commit, review the committed diff again. If the post-commit review finds issues, create a new fix commit and repeat post-commit review until no issues remain.
14. Only after post-commit review finds no issues may the phase todo be marked `completed` and the next phase begin.

## Specialized Review Prompts

When execution reaches the pre-commit review gate, launch five subagents with self-contained packets. Each packet must include the user goal, current phase, relevant diff, tests run, browser validation evidence if frontend changed, and repository constraints.

Use these reviewer scopes:

- Security reviewer: auth bypass, path traversal, HTML injection, XSS, sensitive data leaks, permission boundaries, unsafe logging, credential reads.
- Logic bug reviewer: state flow, error branches, frontend/backend contract, edge cases, data compatibility, regressions.
- Test coverage reviewer: missing unit/integration/frontend/browser coverage, untested high-risk paths, insufficient assertions.
- Maintainability reviewer: module boundaries, naming, duplication, abstraction level, configuration design, documentation, future extension cost.
- Performance/concurrency reviewer: large file reads, memory pressure, caching, concurrent requests, locks, rate limits, timeouts, N+1 calls.

Each reviewer must return either:

- `未发现问题`
- or findings ordered by severity, with concrete file/path references and recommended fixes.

## Browser Screenshot Validation

For any frontend change, the roadmap must require:

- real browser validation, not user screenshots or verbal confirmation
- at least one browser screenshot of the changed page
- desktop viewport coverage
- a smaller/narrow viewport when layout, sidebar, table, card, iframe, modal, or navigation changes
- checks for navigation highlight, route transitions, loading state, empty state, error state, button usability, text truncation, overlap, scroll behavior, and iframe sizing when relevant
- no commit until screenshot validation passes

## Standard Template

Use this structure and replace bracketed placeholders with the user's domain-specific content.

```markdown
---
name: [short-kebab-name]
overview: [one-sentence overview]
todos:
  - id: phase-1-[slug]
    content: Phase 1：[phase title]
    status: pending
  - id: phase-2-[slug]
    content: Phase 2：[phase title]
    status: pending
    depends_on:
      - phase-1-[slug]
isProject: false
---

# [Roadmap Title]

## /goal 提示词：[Goal Name]

将下面整份 Markdown 作为 `/goal` 输入使用。[one-paragraph goal summary]

## 背景

[verified current state and why this work exists]

## 核心目标

- [goal]

## 非目标

- [non-goal]

## 总体原则

- [domain principle]
- 每个阶段必须小步提交；一个阶段一个 commit，不混入无关改动。
- 涉及代码 symbol 编辑前按仓库规则执行影响分析；提交前执行变更检测。

## 强制阶段门禁

[insert Required Global Gates]

## Review 门禁细则

[insert Specialized Review Prompts]

## 浏览器截图验收门禁

[insert Browser Screenshot Validation]

## Phase 1：[title]

### Codex Goal

[phase goal]

### Scope

- [in scope]

### Acceptance Criteria

- [criterion]
- 完成自动化测试和必要浏览器截图验收后，执行 5 个提交前专项 review；全部无问题后才允许提交。
- 提交后再次 review 最新 commit；review 无问题后才允许进入下一阶段。

## 总体最终验收标准

- [final criterion]

## 当前执行建议

下一步只细化 Phase 1 的执行 plan。后续 Phase 保持目标和验收级别，等前一阶段完成、提交后 review 无问题，再基于实际代码状态滚动细化下一阶段。
```

## Authoring Rules

- Use Simplified Chinese for repository-facing plans.
- Keep phase scope narrow and independently committable.
- Put high-risk requirements in global gates, not only in individual phases.
- Do not ask questions inside the plan; ask before creating the plan.
- Do not preserve compatibility with unshipped branch work unless it protects persisted data or public interfaces.
- Do not include secrets, credential paths beyond safe names, or generated artifacts as commit targets.
- If the final phase marks the last todo completed, require `bash scripts/archive-completed-plans.sh` before the final commit when repository rules require plan archiving.
