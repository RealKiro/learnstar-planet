<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>班宠星球 - 班级宠物积分管理系统</title>
    
    <!-- PWA Meta -->
    <meta name="theme-color" content="#6366f1">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
    <meta name="application-name" content="班宠星球">
    <link rel="manifest" href="/manifest.json">
    
    <!-- Livewire Styles -->
    @livewireStyles
    
    <!-- FluxUI -->
    @fluxStyles
    
    <!-- App Styles -->
    <link rel="stylesheet" href="{{ asset('css/app.css') }}">
    
    <!-- Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link href="https://fonts.googleapis.com/css2?family=Geist:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    
    <style>
        /* 全免费标识 - 始终显示 */
        .free-badge {
            background: linear-gradient(135deg, #10b981, #34d399);
            color: white;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            letter-spacing: 0.5px;
        }
        
        /* 班宠星球主题色 */
        :root {
            --bc-primary: #6366f1;
            --bc-primary-light: #818cf8;
            --bc-primary-dark: #4f46e5;
            --bc-secondary: #f59e0b;
            --bc-success: #10b981;
            --bc-danger: #ef4444;
            --bc-info: #3b82f6;
            --bc-bg: #f8fafc;
            --bc-card: #ffffff;
            --bc-text: #1e293b;
            --bc-text-light: #64748b;
        }
        
        [data-theme="dark"] {
            --bc-bg: #0f172a;
            --bc-card: #1e293b;
            --bc-text: #f1f5f9;
            --bc-text-light: #94a3b8;
        }
        
        body {
            font-family: 'Geist', -apple-system, sans-serif;
            background: var(--bc-bg);
            color: var(--bc-text);
            transition: background 0.3s, color 0.3s;
        }
        
        .bc-card {
            background: var(--bc-card);
            border: 1px solid rgba(0,0,0,0.06);
            border-radius: 16px;
            padding: 24px;
            transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.3s;
        }
        
        .bc-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 40px rgba(0,0,0,0.08);
        }
        
        .bc-gradient-text {
            background: linear-gradient(135deg, var(--bc-primary), var(--bc-secondary));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        .bc-glass {
            background: rgba(255,255,255,0.05);
            backdrop-filter: blur(30px) saturate(200%);
            border: 1px solid rgba(255,255,255,0.1);
            border-radius: 20px;
        }
        
        .bc-magnetic {
            transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
        }
        .bc-magnetic:hover {
            transform: scale(1.05) translateY(-2px);
        }
        
        /* 宠物进化路径样式 */
        .pet-evolution-path {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 16px 0;
        }
        .pet-stage {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 4px;
        }
        .pet-emoji {
            font-size: 32px;
            transition: transform 0.3s;
        }
        .pet-emoji:hover {
            transform: scale(1.2);
        }
        .pet-stage.current .pet-emoji {
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.1); }
        }
    </style>
</head>
<body>
    <!-- 导航栏 -->
    <nav class="bc-glass sticky top-0 z-50 px-6 py-4">
        <div class="max-w-7xl mx-auto flex items-center justify-between">
            <div class="flex items-center gap-3">
                <span class="text-2xl">🪐</span>
                <h1 class="text-xl font-bold bc-gradient-text">班宠星球</h1>
                <span class="free-badge">全免费</span>
            </div>
            
            <div class="flex items-center gap-4">
                <!-- 主题切换 -->
                <flux:button.icon variant="ghost" onclick="toggleTheme()">
                    <flux:icon.sun class="hidden dark:block" />
                    <flux:icon.moon class="block dark:hidden" />
                </flux:button.icon>
                
                @auth
                <!-- 用户菜单 -->
                <flux:dropdown>
                    <flux:button variant="ghost" icon:trailing="chevron-down">
                        {{ auth()->user()->name }}
                    </flux:button>
                    <flux:menu>
                        <flux:menu.item href="/dashboard">数据看板</flux:menu.item>
                        <flux:menu.item href="/settings">设置</flux:menu.item>
                        <flux:menu.separator />
                        <flux:menu.item wire:click="logout" variant="danger">退出登录</flux:menu.item>
                    </flux:menu>
                </flux:dropdown>
                @else
                <flux:button href="/login" variant="primary">登录</flux:button>
                <flux:button href="/register" variant="ghost">注册</flux:button>
                @endauth
            </div>
        </div>
    </nav>

    <!-- 主内容 -->
    <main class="max-w-7xl mx-auto px-6 py-8">
        {{ $slot }}
    </main>

    <!-- 底部 -->
    <footer class="border-t py-8 mt-16">
        <div class="max-w-7xl mx-auto px-6 text-center text-sm opacity-60">
            <p>班宠星球 - 让每个班级都充满乐趣 | 全免费开源</p>
            <p class="mt-2">支持 Windows / Mac / Android / 微信小程序</p>
        </div>
    </footer>

    <!-- Livewire Scripts -->
    @livewireScripts
    @fluxScripts
    
    <!-- PWA Service Worker -->
    <script>
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.register('/sw.js');
        }
        
        function toggleTheme() {
            const current = document.documentElement.getAttribute('data-theme') || 'light';
            const next = current === 'light' ? 'dark' : 'light';
            document.documentElement.setAttribute('data-theme', next);
            localStorage.setItem('theme', next);
        }
        
        // 初始化主题
        const savedTheme = localStorage.getItem('theme') || 'system';
        if (savedTheme === 'system') {
            document.documentElement.setAttribute('data-theme', 
                window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
            );
        } else {
            document.documentElement.setAttribute('data-theme', savedTheme);
        }
    </script>
</body>
</html>
