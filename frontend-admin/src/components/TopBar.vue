<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const menuOpen = ref(false)

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <header class="masthead">
    <div class="brand">
      <span class="brand-title">Zimbra Administration</span>
      <span class="brand-sub">A&nbsp;SYNACOR&nbsp;PRODUCT</span>
    </div>

    <div class="search">
      <span class="search-type" aria-hidden="true"></span>
      <input type="text" placeholder="" aria-label="Search" />
      <button class="search-go" aria-label="Search">&#128269;</button>
    </div>

    <div class="masthead-actions">
      <button class="mh-icon" aria-label="Refresh">&#8635;</button>
      <div class="mh-user" @click="menuOpen = !menuOpen">
        <span>{{ auth.username || 'admin' }}</span>
        <span class="caret">&#9662;</span>
        <ul v-if="menuOpen" class="mh-menu" @click.stop>
          <li @click="logout">Sign Out</li>
        </ul>
      </div>
      <div class="mh-help">
        <span>Help</span><span class="caret">&#9662;</span>
      </div>
    </div>
  </header>
</template>

<style scoped>
.masthead {
  display: flex;
  align-items: center;
  height: 44px;
  padding: 0 12px;
  color: var(--masthead-txt);
  background: linear-gradient(to bottom, var(--masthead-top), var(--masthead-bottom));
  gap: 16px;
}
.brand {
  display: flex;
  flex-direction: column;
  line-height: 1.05;
  min-width: 220px;
}
.brand-title {
  font-size: 19px;
  font-weight: 400;
  letter-spacing: 0.2px;
}
.brand-sub {
  font-size: 8px;
  letter-spacing: 1px;
  opacity: 0.85;
}
.search {
  display: flex;
  align-items: center;
  flex: 1;
  max-width: 560px;
  height: 24px;
  margin: 0 auto;
  background: #fff;
  border-radius: var(--radius);
  padding: 0 4px;
}
.search-type {
  width: 24px;
  height: 18px;
  border-right: 1px solid #ccc;
}
.search input {
  flex: 1;
  border: none;
  outline: none;
  height: 22px;
  font-size: var(--fs);
  padding: 0 6px;
  color: var(--txt);
}
.search-go {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 13px;
}
.masthead-actions {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: var(--fs);
}
.mh-icon {
  background: none;
  border: none;
  color: #fff;
  font-size: 15px;
  cursor: pointer;
}
.mh-user,
.mh-help {
  position: relative;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}
.caret {
  font-size: 9px;
}
.mh-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin: 4px 0 0;
  padding: 4px 0;
  list-style: none;
  background: #fff;
  color: var(--txt);
  border: 1px solid var(--panel-border);
  border-radius: var(--radius);
  min-width: 120px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  z-index: 50;
}
.mh-menu li {
  padding: 5px 14px;
}
.mh-menu li:hover {
  background: var(--sel);
}
</style>
