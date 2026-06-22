// ============================================================
// 班宠星球 - PWA Service Worker
// 支持 Windows/Mac/Android 离线安装
// ============================================================

const CACHE_NAME = 'bancxq-planet-v1';
const STATIC_ASSETS = [
  '/',
  '/manifest.json',
  '/css/app.css',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png'
];

// 安装：预缓存静态资源
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(STATIC_ASSETS))
      .then(() => self.skipWaiting())
  );
});

// 激活：清理旧缓存
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys.filter((key) => key !== CACHE_NAME)
          .map((key) => caches.delete(key))
      );
    }).then(() => self.clients.claim())
  );
});

// 请求拦截：网络优先，缓存后备
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // API 请求：仅网络，不缓存
  if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/livewire/')) {
    return event.respondWith(
      fetch(request).catch(() => new Response(JSON.stringify({ error: 'offline' }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' }
      }))
    );
  }

  // SSE 请求：不拦截
  if (url.pathname.startsWith('/api/events')) {
    return;
  }

  // 静态资源：缓存优先
  if (request.method === 'GET') {
    event.respondWith(
      caches.match(request).then((cached) => {
        if (cached) return cached;
        return fetch(request).then((response) => {
          // 只缓存成功响应
          if (response.status === 200) {
            const clone = response.clone();
            caches.open(CACHE_NAME).then((cache) => cache.put(request, clone));
          }
          return response;
        });
      })
    );
  }
});

// 推送通知
self.addEventListener('push', (event) => {
  const data = event.data ? event.data.json() : {};
  
  const title = data.title || '班宠星球通知';
  const options = {
    body: data.body || '你有新的通知消息',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/icon-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      url: data.url || '/'
    },
    actions: [
      { action: 'open', title: '查看详情' },
      { action: 'close', title: '忽略' }
    ]
  };

  event.waitUntil(
    self.registration.showNotification(title, options)
  );
});

// 通知点击
self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  
  if (event.action === 'open') {
    event.waitUntil(
      clients.openWindow(event.notification.data.url || '/')
    );
  }
});

// 后台同步
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-scores') {
    event.waitUntil(syncPendingScores());
  }
});

async function syncPendingScores() {
  // 同步离线期间的积分操作
  const pending = await getPendingOperations();
  for (const op of pending) {
    try {
      await fetch('/api/teacher/scores/quick', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${op.token}`
        },
        body: JSON.stringify(op.data)
      });
      await removePendingOperation(op.id);
    } catch (err) {
      console.error('Sync failed:', err);
    }
  }
}
