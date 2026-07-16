# Product Goals (No-PRD)

**Project**: gzh-cli-mcp-plugin (`mcp-plugin` binary)
**Doc Type**: Goals + Constraints + Quality Gates
**Status**: Active — 방향 확정 전, 최소 계약
**Last Updated**: 2026-07-16

______________________________________________________________________

## Product Intent

gzh-cli-mcp-plugin is an **MCP configuration manager for Claude Code**. It:

- edits Claude Code's config files (`~/.claude.json`, `~/.claude/settings.json`)
  to list/install/remove MCP servers and toggle plugins,
- discovers MCP packages on the npm registry (read-only),
- and **never speaks the MCP protocol** — Claude Code owns server lifecycle;
  this tool wraps its configuration (SOUL 신념 1).

This is a feature-library project — a single PRODUCT.md is sufficient. It
replaces a PRD.

| 제공하는 것 (Is)                              | 되지 않을 것 (Is Not)                       |
| --------------------------------------------- | ------------------------------------------- |
| Claude Code MCP 설정 파일 읽기/쓰기           | MCP 서버·클라이언트·트랜스포트 구현         |
| MCP 서버 등록·제거·플러그인 토글              | MCP 서버 라이프사이클 관리                  |
| npm 레지스트리 검색·조회 (읽기 전용)          | 범용 npm 클라이언트·패키지 설치             |
| 설정 export/import/validate                   | Claude Code 내부 수정                       |

______________________________________________________________________

## Goals (Measurable Targets)

G1. **Config-manager scope**

- Target: Claude Code 설정 파일 3종(`~/.claude.json`,
  `~/.claude/settings.json`, `plugins/cache/*/.mcp.json`)에 대한 읽기·쓰기를
  11개 명령으로 제공 — 현재 충족 (build·vet·lint 모두 clean)

G2. **Test coverage on the shipped path**

- Target: 실사용 경로 `pkg/config` 커버리지 >= 50%
- 현재 **15.3%**. 전체 10.9%이며, 높은 수치(usecase 85.4% · claudeconfig 73.3%)는
  전부 **미연결 코드**에 붙어 있어 합산 커버리지가 실제 가치를 과대평가한다

G3. **Zero orphaned packages**

- Target: import되지 않는 패키지 0개
- 현재 **4개** — `application/usecase/plugin`, `infrastructure/repository`,
  `domain/mcp`, `infrastructure/claudeconfig`. 11개 명령이 모두 `pkg/config`를
  직접 호출하며 Ports&Adapters 골격은 어디에도 연결돼 있지 않다.
  **연결하거나 삭제해야 한다 — 둘 중 하나를 고르는 것이 본 리포 최우선 결정이다**

G4. **Command coherence**

- Target: `list`가 나열하는 대상과 `enable`/`disable`이 조작하는 대상이 일치
- 현재 **불일치** — `list`는 MCP *서버*를, `enable`/`disable`은 *플러그인*을
  다루는데 도움말은 양쪽 모두 "`list`로 플러그인을 확인하라"고 안내한다

G5. **Baseline alignment**

- Target: `.golangci.yml`(v2) 도입 + `go.mod` go directive를 devbox 툴체인(1.26)에
  정렬 — 현재 미충족 (golangci 설정 없음, go 1.24)

______________________________________________________________________

## Non-Goals (Explicitly Out of Scope)

- No MCP 서버·클라이언트·트랜스포트 구현 — MCP SDK를 의존하지 않는다
- No MCP 서버 라이프사이클 관리 — Claude Code가 소유한다
- No npm·Python 패키지 설치/제거 — 설정 항목만 다룬다
- No Claude Code 내부 수정
- No 범용 npm 클라이언트 — npm 접근은 검색·조회(읽기 전용)뿐
- No CGO 의존

______________________________________________________________________

## Guardrails and Technical Constraints

**Architecture**

- **현재 실제 구현은 CLI → `pkg/config` 직접 경로다.** 리포의 CLAUDE.md는
  "CLI → infrastructure 직접 호출 금지"를 규정하지만 코드는 그 규칙을 지키지 않는다.
  규칙을 코드에 맞추거나 코드를 규칙에 맞추기 전까지 **새 계층을 추가하지 않는다**

**Dependency Boundaries**

- `gzh-cli-core`만 의존 가능; 다른 feature 라이브러리 의존 금지 (GUIDELINES §2)
- 현재 직접 의존은 cobra 뿐 — CLAUDE.md의 core 사용 안내는 코드와 맞지 않는다

**Compatibility**

- Go 1.24 (`go.mod`) — devbox 툴체인 1.26 대비 하위. GUIDELINES §1에 따라 릴리스
  시점에 상향 정렬한다

**Safety**

- 사용자 홈의 Claude Code 설정을 직접 수정한다 — 모든 변경은 재시작 안내를 출력한다
- 서버 정보 출력 시 auth/token/key 헤더 값을 마스킹한다
- 설정 파일 백업·롤백은 없다 — 파괴적 범위를 넓히기 전에 도입해야 한다

**Baseline (진행 중)**

- `.golangci.yml` 미보유 (GUIDELINES §4 격차). 현재 `make lint`의 "0 issues"는
  golangci-lint **기본 린터**만 돈 결과이며, `.make/vars.mk`(v1.62.2)와
  CI(`latest`)가 서로 다른 버전을 가리킨다 — 공유 룰셋이 없다

______________________________________________________________________

## Quality Gates (Release Readiness)

**Build and Lint**

- `make validate` (fmt + vet + lint + test) pass with no warnings

**Testing**

- `pkg/config` 커버리지 >= 50% (G2)

**Docs**

- `CLAUDE.md`가 실제 명령·구조와 일치한다 (현재 미충족: "Initial stage - project
  structure only", "MVP Scope: list/status/enable/disable" 서술이 stale —
  이미 11개 명령이 쓰기 연산까지 제공한다)

______________________________________________________________________

## Decision Rules

- **MCP 프로토콜 자체 구현은 SOUL 게이트 1(재발명 금지)에서 거절된다** — Claude
  Code를 감쌀 뿐이다
- 미연결 계층(G3)에 코드를 추가하지 않는다 — 먼저 연결 여부를 결정한다
- 새 기능은 SOUL.md 4-게이트(틈 · 라이브러리 · 대량/전환 · 날카로움)를 통과해야 한다
- Guardrails 위반은 문서화된 예외를 요구한다
- Quality Gates 미충족 시 릴리스는 차단된다

______________________________________________________________________

**End of Document**
