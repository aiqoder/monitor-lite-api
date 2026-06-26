import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../layout/index.vue'
import RuleEnginaView from '@/views/ruleEngina/index.vue'
import tvlistView from '@/views/tvlist/index.vue'
import settingView from '@/views/setting/index.vue'
import LoginView from '@/views/Login.vue'
import SearchView from '@/views/search/index.vue'
import SubscriberView from '@/views/subscriber/index.vue'
import epgView from '@/views/epg/index.vue'
import onlineView from '@/views/online/index.vue'

const router = createRouter({
  history: createWebHistory("/admin"),
  routes: [
    {
      path: '/',
      name: 'layout',
      component: Layout,
      redirect: {
        path: "tvlist"
      },
      children: [
        {
          path: "rule",
          component: RuleEnginaView,
        },
        {
          path: "search",
          component: SearchView,
        },
        {
          path: "tvlist",
          component: tvlistView,
        },
        {
          path: "subscriber",
          component: SubscriberView, // 订阅池
        },
        {
          path: "setting",
          component: settingView,
        },
        {
          path: "epg",
          component: epgView,
        },
        {
          path: "online",
          component: onlineView,
        },
      ]
    },
    {
      path: "/login",
      component: LoginView,
    },
  ]
})

router.beforeEach((to, from) => {
  if (to.path === '/login') {
    return true
  }

  const pwd = sessionStorage.getItem("auth")

  if(!pwd) {
    return "/login"
  }else {
    return true
  }
})

export default router
