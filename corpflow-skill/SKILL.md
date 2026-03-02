# CorpFlow Skill

企业级 AI 协作平台 - 多智能体团队

## 概述

CorpFlow 是一个 AI 驱动的公司治理平台，提供三层智能体协作：
- **CEO** - 战略决策、目标制定
- **Manager** - 任务分配、进度跟踪
- **Worker** - 执行任务、汇报结果

## 核心特性

### 🧠 动态模型选择

根据任务类型和复杂度自动选择最合适的模型：

| 模型 | 用途 | 特点 |
|------|------|------|
| **GLM-5** | 战略规划、复杂决策 | 长上下文、深度推理 |
| **MiniMax-M2.5** | 任务分解、进度评估 | 逻辑清晰、分析能力强 |
| **Kimi K2.5** | 简单执行、信息检索 | 快速响应、高效处理 |

### 📊 绩效模型切换

根据 Agent 历史表现自动调整模型：

- 绩效 > 90 分 → 升级模型
- 绩效 70-90 分 → 保持现状
- 绩效 < 70 分 → 降级稳定模型

## 安装

### 方式 1: 克隆到 skills 目录

```bash
# 复制 skill 文件到 OpenClaw
cp -r corpflow-skill/* ~/.openclaw/workspace/skills/corpflow/
```

### 方式 2: Git 子模块

```bash
git submodule add https://github.com/gotonote/corpflow.git skills/corpflow
```

## 使用

在飞书与 OpenClaw 对话时，可直接调用：

```
@OpenClaw 帮我制定本季度 OKR
@OpenClaw 分析当前项目进度
@OpenClaw 分配任务给团队
```

AI 会自动选择合适的模型执行任务。

## 定时任务

已配置自动任务：
- 工作日 9:00 - OKR 进度提醒
- 周五 18:00 - 周总结

---

*Version: 1.1 - 支持动态模型选择*
