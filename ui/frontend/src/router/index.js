import { createRouter, createWebHistory } from "vue-router";
import axios from "axios";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("../views/HomeView.vue"),
    },
    {
      path: "/user/login",
      name: "login",
      component: () => import("../views/LoginView.vue"),
    },
    {
      path: "/user/signup",
      name: "signup",
      component: () => import("../views/SignUpView.vue"),
    },
    {
      path: "/dashboard",
      name: "dashboard",
      component: () => import("../views/DashboardView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/privacy-policy",
      name: "privacy-policy",
      component: () => import("../views/PrivacyPolicyView.vue"),
    },
    {
      path: "/:user",
      name: "superbchat",
      component: () => import("../views/SuperbchatView.vue"),
    },
    {
      path: "/404",
      name: "404",
      component: () => import("../views/404View.vue"),
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    const token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await axios.post("/verify-tkn", { token });
        if (response.statusText == "OK") {
          next(); 
        } else {
          next("/user/login"); 
        }
      } catch (error) {
        console.error("Error verificando el token:", error);
        next("/user/login"); 
      }
    } else {
      next("/user/login"); 
    }
  } else {
    next(); 
  }
});

export default router;
