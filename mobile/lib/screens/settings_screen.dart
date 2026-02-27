import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  final _serverController = TextEditingController(text: 'http://localhost:8080');
  bool _darkMode = false;
  bool _multiModelVoting = true;
  String _defaultModel = 'gpt-4';
  String _votingMethod = 'comprehensive';

  @override
  void initState() {
    super.initState();
    _loadSettings();
  }

  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      _serverController.text = prefs.getString('server_url') ?? 'http://localhost:8080';
      _darkMode = prefs.getBool('dark_mode') ?? false;
      _multiModelVoting = prefs.getBool('multi_model_voting') ?? true;
      _defaultModel = prefs.getString('default_model') ?? 'gpt-4';
      _votingMethod = prefs.getString('voting_method') ?? 'comprehensive';
    });
  }

  Future<void> _saveSettings() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('server_url', _serverController.text);
    await prefs.setBool('dark_mode', _darkMode);
    await prefs.setBool('multi_model_voting', _multiModelVoting);
    await prefs.setString('default_model', _defaultModel);
    await prefs.setString('voting_method', _votingMethod);
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Settings saved / 设置已保存')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings / 设置'),
      ),
      body: ListView(
        children: [
          // Server Settings / 服务器设置
          _buildSectionHeader('Server / 服务器'),
          ListTile(
            leading: const Icon(Icons.dns),
            title: const Text('Server URL / 服务器地址'),
            subtitle: Text(_serverController.text),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showServerDialog(),
          ),
          
          // Appearance / 外观
          _buildSectionHeader('Appearance / 外观'),
          SwitchListTile(
            secondary: const Icon(Icons.dark_mode),
            title: const Text('Dark Mode'),
            subtitle: const Text('深色模式'),
            value: _darkMode,
            onChanged: (value) {
              setState(() => _darkMode = value);
              _saveSettings();
            },
          ),
          
          // Model Settings / 模型设置
          _buildSectionHeader('AI Models / AI模型'),
          ListTile(
            leading: const Icon(Icons.model_training),
            title: const Text('Default Model / 默认模型'),
            subtitle: Text(_defaultModel),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showModelDialog(),
          ),
          const Divider(),
          
          // Multi-Model Voting / 多模型投票
          _buildSectionHeader('🗳️ Multi-Model Voting / 多模型投票'),
          SwitchListTile(
            secondary: const Icon(Icons.poll),
            title: const Text('Enable Multi-Model Voting'),
            subtitle: const Text('启用多模型投票决策'),
            value: _multiModelVoting,
            onChanged: (value) {
              setState(() => _multiModelVoting = value);
              _saveSettings();
            },
          ),
          if (_multiModelVoting) ...[
            ListTile(
              leading: const Icon(Icons.analytics),
              title: const Text('Voting Method / 投票方式'),
              subtitle: Text(_getVotingMethodName(_votingMethod)),
              trailing: const Icon(Icons.chevron_right),
              onTap: () => _showVotingMethodDialog(),
            ),
            _buildVotingMethodInfo(),
          ],
          
          // Supported Models / 支持的模型
          _buildSectionHeader('📋 Supported Models / 支持的模型'),
          _buildModelList(),
          
          // Channels / 渠道
          _buildSectionHeader('📱 Channels / 渠道'),
          SwitchListTile(
            secondary: const Icon(Icons.chat_bubble),
            title: const Text('Feishu / 飞书'),
            value: true,
            onChanged: (v) {},
          ),
          SwitchListTile(
            secondary: const Icon(Icons.chat_bubble_outline),
            title: const Text('WeChat / 微信'),
            value: true,
            onChanged: (v) {},
          ),
          SwitchListTile(
            secondary: const Icon(Icons.send),
            title: const Text('Telegram'),
            value: false,
            onChanged: (v) {},
          ),
          SwitchListTile(
            secondary: const Icon(Icons.discord),
            title: const Text('Discord'),
            value: false,
            onChanged: (v) {},
          ),
          
          // API Keys / API密钥
          _buildSectionHeader('🔑 API Keys / API密钥'),
          ListTile(
            leading: const Icon(Icons.key),
            title: const Text('Configure API Keys / 配置API密钥'),
            subtitle: const Text('Set environment variables on server'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showApiKeysInfo(),
          ),
          
          // About / 关于
          _buildSectionHeader('About / 关于'),
          const ListTile(
            leading: Icon(Icons.info),
            title: Text('Version'),
            subtitle: Text('1.0.0'),
          ),
          const ListTile(
            leading: Icon(Icons.code),
            title: Text('GitHub'),
            subtitle: Text('github.com/gotonote/corpflow'),
            trailing: Icon(Icons.open_in_new),
          ),
          
          const SizedBox(height: 32),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: ElevatedButton(
              onPressed: _saveSettings,
              child: const Text('Save Settings / 保存设置'),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
      child: Text(
        title,
        style: TextStyle(
          color: Theme.of(context).primaryColor,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Widget _buildVotingMethodInfo() {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Card(
        color: Colors.blue.shade50,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Voting Methods / 投票方式:',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              Text('• Comprehensive (默认): Accuracy + Completeness + Clarity + Creativity'),
              Text('• Cross-evaluation: Models evaluate each other'),
              Text('• Length: Simply by response length'),
              const SizedBox(height: 8),
              Text(
                '评估维度 / Evaluation Dimensions:',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 4),
              Text('• Accuracy (准确性) - 30%'),
              Text('• Completeness (完整性) - 30%'),
              Text('• Clarity (清晰度) - 20%'),
              Text('• Creativity (创造性) - 20%'),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildModelList() {
    final models = [
      {'name': 'GPT-4', 'provider': 'OpenAI', 'key': 'OPENAI_API_KEY'},
      {'name': 'Claude 3', 'provider': 'Anthropic', 'key': 'ANTHROPIC_API_KEY'},
      {'name': 'GLM-4', 'provider': 'Zhipu (智谱)', 'key': 'ZHIPU_API_KEY'},
      {'name': 'Kimi', 'provider': 'Moonshot (月之暗面)', 'key': 'KIMI_API_KEY'},
      {'name': 'Qwen', 'provider': 'Alibaba (通义千问)', 'key': 'DASHSCOPE_API_KEY'},
      {'name': 'DeepSeek', 'provider': 'DeepSeek', 'key': 'DEEPSEEK_API_KEY'},
      {'name': 'MiniMax', 'provider': 'MiniMax', 'key': 'MINIMAX_API_KEY'},
    ];

    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        children: models.map((model) => ListTile(
          leading: const Icon(Icons.psychology),
          title: Text(model['name']!),
          subtitle: Text(model['provider']!),
          trailing: IconButton(
            icon: const Icon(Icons.info_outline),
            onPressed: () {},
          ),
        )).toList(),
      ),
    );
  }

  String _getVotingMethodName(String method) {
    switch (method) {
      case 'comprehensive':
        return 'Comprehensive (综合评分)';
      case 'cross':
        return 'Cross-evaluation (交叉评估)';
      case 'length':
        return 'By Length (按长度)';
      default:
        return method;
    }
  }

  void _showServerDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Server URL / 服务器地址'),
        content: TextField(
          controller: _serverController,
          decoration: const InputDecoration(
            labelText: 'URL',
            hintText: 'http://localhost:8080',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              _saveSettings();
              Navigator.pop(context);
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  void _showModelDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Select Default Model / 选择默认模型'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            'gpt-4',
            'claude-3-opus',
            'glm-4',
            'moonshot-v1-8k-chat',
            'qwen-turbo',
            'deepseek-chat',
            'abab6.5s-chat',
          ].map((model) => RadioListTile<String>(
            title: Text(model),
            value: model,
            groupValue: _defaultModel,
            onChanged: (v) {
              setState(() => _defaultModel = v!);
              Navigator.pop(context);
              _saveSettings();
            },
          )).toList(),
        ),
      ),
    );
  }

  void _showVotingMethodDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Voting Method / 投票方式'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<String>(
              title: const Text('Comprehensive (综合评分)'),
              subtitle: const Text('Accuracy/Completeness/Clarity/Creativity'),
              value: 'comprehensive',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
            RadioListTile<String>(
              title: const Text('Cross-evaluation (交叉评估)'),
              subtitle: const Text('Models evaluate each other'),
              value: 'cross',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
            RadioListTile<String>(
              title: const Text('By Length (按长度)'),
              subtitle: const Text('Simple length-based'),
              value: 'length',
              groupValue: _votingMethod,
              onChanged: (v) {
                setState(() => _votingMethod = v!);
                Navigator.pop(context);
                _saveSettings();
              },
            ),
          ],
        ),
      ),
    );
  }

  void _showApiKeysInfo() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('API Keys Configuration / API密钥配置'),
        content: const SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('Set these environment variables on your server:',
                style: TextStyle(fontWeight: FontWeight.bold)),
              SizedBox(height: 12),
              Text('# OpenAI\nexport OPENAI_API_KEY=sk-xxx'),
              SizedBox(height: 8),
              Text('# Anthropic\nexport ANTHROPIC_API_KEY=sk-ant-xxx'),
              SizedBox(height: 8),
              Text('# Zhipu GLM\nexport ZHIPU_API_KEY=xxx'),
              SizedBox(height: 8),
              Text('# Kimi (Moonshot)\nexport KIMI_API_KEY=xxx'),
              SizedBox(height: 8),
              Text('# Qwen (Alibaba)\nexport DASHSCOPE_API_KEY=xxx'),
              SizedBox(height: 8),
              Text('# DeepSeek\nexport DEEPSEEK_API_KEY=xxx'),
              SizedBox(height: 8),
              Text('# MiniMax\nexport MINIMAX_API_KEY=xxx'),
              SizedBox(height: 16),
              Text('For Docker, add to docker-compose.yml environment section.',
                style: TextStyle(color: Colors.grey)),
            ],
          ),
        ),
        actions: [
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _serverController.dispose();
    super.dispose();
  }
}
