import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('CorpFlow'),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Welcome Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.smart_toy, size: 40, color: Color(0xFF667eea)),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Welcome to CorpFlow',
                                style: Theme.of(context).textTheme.headlineSmall,
                              ),
                              Text(
                                'Multi-Agent Collaboration Platform',
                                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                  color: Colors.grey[600],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),
            
            // Quick Actions / 快捷操作
            Text(
              'Quick Actions',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: _QuickActionCard(
                    icon: Icons.chat,
                    title: 'New Chat',
                    subtitle: 'Start conversation',
                    onTap: () {},
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _QuickActionCard(
                    icon: Icons.account_tree,
                    title: 'New Flow',
                    subtitle: 'Create workflow',
                    onTap: () {},
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            
            // Statistics
            Text(
              'Statistics',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(child: _StatCard(title: 'Agents', value: '5', icon: Icons.smart_toy)),
                const SizedBox(width: 12),
                Expanded(child: _StatCard(title: 'Flows', value: '12', icon: Icons.account_tree)),
                const SizedBox(width: 12),
                Expanded(child: _StatCard(title: 'Messages', value: '1.2k', icon: Icons.chat)),
              ],
            ),
            const SizedBox(height: 20),
            
            // Features Guide / 功能指南
            Text(
              'Features Guide',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            
            // Chat Feature
            _FeatureCard(
              icon: Icons.chat,
              title: '💬 Chat / 对话',
              description: 'Chat with AI agents. Support multi-channel (Feishu, WeChat, Telegram).',
              descriptionZh: '与AI智能体对话。支持多渠道接入（飞书、微信、Telegram）。',
              steps: [
                '1. Tap "New Chat" to start / 点击"新建对话"开始',
                '2. Type your message / 输入你的消息',
                '3. AI responds instantly / AI即时回复',
              ],
            ),
            const SizedBox(height: 12),
            
            // Flow Feature
            _FeatureCard(
              icon: Icons.account_tree,
              title: '🔀 Flow / 流程编排',
              description: 'Visual workflow editor. Drag-and-drop nodes to create automation.',
              descriptionZh: '可视化流程编辑器。拖拽节点创建自动化工作流。',
              steps: [
                '1. Tap "New Flow" / 点击"新建流程"',
                '2. Add nodes (Trigger/Agent/Tool) / 添加节点',
                '3. Connect nodes / 连接节点',
                '4. Save and execute / 保存并执行',
              ],
            ),
            const SizedBox(height: 12),
            
            // Agent Feature
            _FeatureCard(
              icon: Icons.smart_toy,
              title: '🤖 Agents / 智能体',
              description: 'Create and manage AI agents. Configure models and tools.',
              descriptionZh: '体。配置模型创建和管理AI智能和工具。',
              steps: [
                '1. Go to Agents tab / 进入智能体标签页',
                '2. Tap "+" to create / 点击"+"创建',
                '3. Select AI model (GPT-4/Claude/GLM-4/Kimi/Qwen/DeepSeek)',
                '4. Configure tools / 配置工具',
              ],
            ),
            const SizedBox(height: 12),
            
            // Multi-Model Voting
            _FeatureCard(
              icon: Icons.poll,
              title: '🗳️ Multi-Model Voting / 多模型投票',
              description: 'Let multiple AI models discuss and vote on best decision.',
              descriptionZh: '让多个AI模型讨论并投票选择最佳决策。',
              steps: [
                '1. Enable voting in Settings / 在设置中启用投票',
                '2. Select models / 选择模型',
                '3. System evaluates: Accuracy / Completeness / Clarity / Creativity',
                '4. Auto-select best response / 自动选择最佳回复',
              ],
            ),
            const SizedBox(height: 20),
            
            // Recent Activity
            Text(
              'Recent Activity',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Card(
              child: ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: 3,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) => ListTile(
                  leading: CircleAvatar(
                    backgroundColor: Color(0xFF667eea).withOpacity(0.1),
                    child: const Icon(Icons.chat, color: Color(0xFF667eea)),
                  ),
                  title: Text('Conversation ${index + 1}'),
                  subtitle: Text('Last message: 5 min ago'),
                  trailing: const Icon(Icons.chevron_right),
                ),
              ),
            ),
            const SizedBox(height: 20),
            
            // Links / 链接
            Text(
              'Links / 链接',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Card(
              child: Column(
                children: [
                  ListTile(
                    leading: const Icon(Icons.language),
                    title: const Text('Documentation'),
                    subtitle: const Text('View docs / 查看文档'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                  const Divider(height: 1),
                  ListTile(
                    leading: const Icon(Icons.code),
                    title: const Text('GitHub'),
                    subtitle: const Text('Source code / 源代码'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                  const Divider(height: 1),
                  ListTile(
                    leading: const Icon(Icons.forum),
                    title: const Text('Discord'),
                    subtitle: const Text('Community / 社区'),
                    trailing: const Icon(Icons.open_in_new),
                    onTap: () {},
                  ),
                ],
              ),
            ),
            
            // Model Support Info
            const SizedBox(height: 20),
            Text(
              'Supported AI Models / 支持的AI模型',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: [
                    _ModelChip(name: 'GPT-4', provider: 'OpenAI'),
                    _ModelChip(name: 'Claude 3', provider: 'Anthropic'),
                    _ModelChip(name: 'GLM-4', provider: 'Zhipu'),
                    _ModelChip(name: 'Kimi', provider: 'Moonshot'),
                    _ModelChip(name: 'Qwen', provider: 'Alibaba'),
                    _ModelChip(name: 'DeepSeek', provider: 'DeepSeek'),
                    _ModelChip(name: 'MiniMax', provider: 'MiniMax'),
                  ],
                ),
              ),
            ),
            
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }
}

class _QuickActionCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;

  const _QuickActionCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              Icon(icon, size: 32, color: Color(0xFF667eea)),
              const SizedBox(height: 8),
              Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
              Text(subtitle, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;

  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Icon(icon, color: Color(0xFF667eea)),
            const SizedBox(height: 4),
            Text(value, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            Text(title, style: TextStyle(fontSize: 12, color: Colors.grey[600])),
          ],
        ),
      ),
    );
  }
}

class _FeatureCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String description;
  final String descriptionZh;
  final List<String> steps;

  const _FeatureCard({
    required this.icon,
    required this.title,
    required this.description,
    required this.descriptionZh,
    required this.steps,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ExpansionTile(
        leading: Icon(icon, color: Color(0xFF667eea)),
        title: Text(title),
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(description),
                const SizedBox(height: 4),
                Text(
                  descriptionZh,
                  style: TextStyle(fontSize: 12, color: Colors.grey[600]),
                ),
                const SizedBox(height: 12),
                const Text(
                  'How to use:',
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                ...steps.map((step) => Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Text(step, style: TextStyle(fontSize: 13, color: Colors.grey[700])),
                )),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ModelChip extends StatelessWidget {
  final String name;
  final String provider;

  const _ModelChip({required this.name, required this.provider});

  @override
  Widget build(BuildContext context) {
    return Chip(
      avatar: const Icon(Icons.psychology, size: 16),
      label: Text('$name ($provider)'),
      backgroundColor: Color(0xFF667eea).withOpacity(0.1),
    );
  }
}
