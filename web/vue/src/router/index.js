import { createWebHashHistory, createRouter } from 'vue-router'

import store from '../store/index'

export const routes = [
  {
    path: '/:pathMatch(.*)*',
    component: () => import('../pages/notFound.vue'),
    meta: {
      noLogin: true,
      noNeedAdmin: true
    }
  },
  {
    path: '/',
    redirect: '/task'
  },
  {
    path: '/install',
    name: 'install',
    component: () => import('../pages/install/index.vue'),
    meta: {
      noLogin: true,
      noNeedAdmin: true
    }
  },
  {
    path: '/task',
    redirect: '/task/list'
  },
  {
    path: '/task/list',
    name: 'task--list',
    component: () => import('../pages/task/list.vue'),
    meta: {
      noNeedAdmin: true
    }
  },
  {
    path: '/task/create',
    name: 'task-create',
    component: () => import('../pages/task/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/task/edit/:id',
    name: 'task-edit',
    component: () => import('../pages/task/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/task/log',
    name: 'task--log',
    component: () => import('../pages/taskLog/list.vue'),
    meta: {
      noNeedAdmin: true
    }
  },
  {
    path: '/task/dependency',
    name: 'task--dependency',
    component: () => import('../pages/task/dependency.vue'),
    meta: {
      noNeedAdmin: true
    }
  },
  {
    path: '/host',
    redirect: '/host/list'
  },
  {
    path: '/host/list',
    name: 'host--list',
    component: () => import('../pages/host/list.vue'),
    meta: {
      noNeedAdmin: true
    }
  },
  {
    path: '/host/create',
    name: 'host-create',
    component: () => import('../pages/host/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/host/edit/:id',
    name: 'host--edit--',
    component: () => import('../pages/host/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/user',
    redirect: '/user/list'
  },
  {
    path: '/user/list',
    name: 'user--list',
    component: () => import('../pages/user/list.vue')
  },
  {
    path: '/user/create',
    name: 'user-create',
    component: () => import('../pages/user/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/user/edit/:id',
    name: 'user-edit',
    component: () => import('../pages/user/edit.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/user/login',
    name: 'user-login',
    component: () => import('../pages/user/login.vue'),
    meta: {
      noLogin: true,
      inSidebar: false
    }
  },
  {
    path: '/user/edit-password/:id',
    name: 'user-edit-password',
    component: () => import('../pages/user/editPassword.vue'),
    meta: {
      inSidebar: false
    }
  },
  {
    path: '/user/edit-my-password',
    name: 'user-edit-my-password',
    component: () => import('../pages/user/editMyPassword.vue'),
    meta: {
      noNeedAdmin: true,
      inSidebar: false
    }
  },
  {
    path: '/system',
    redirect: '/system/notification'
  },
  {
    path: '/system/notification',
    name: 'system--notification',
    component: () => import('../pages/system/notification/email.vue')
  },
  {
    path: '/system/notification/email',
    name: 'system-notification-email',
    component: () => import('../pages/system/notification/email.vue')
  },
  {
    path: '/system/notification/slack',
    name: 'system-notification-slack',
    component: () => import('../pages/system/notification/slack.vue')
  },
  {
    path: '/system/notification/webhook',
    name: 'system-notification-webhook',
    component: () => import('../pages/system/notification/webhook.vue')
  },
  {
    path: '/system/login-log',
    name: 'system--login-log',
    component: () => import('../pages/system/loginLog.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
router.beforeEach((to, from, next) => {
  if (to.meta.noLogin) {
    next()
    return
  }
  if (store.getters.user.token) {
    if ((store.getters.user.isAdmin || to.meta.noNeedAdmin)) {
      next()
      return
    }
    if (!store.getters.user.isAdmin) {
      next(
        {
          path: '/404.html'
        }
      )
      return
    }
  }
  next({
    path: '/user/login',
    query: {redirect: to.fullPath}
  })
})

export default router
