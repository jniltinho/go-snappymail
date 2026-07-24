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
    <div class="hist" aria-hidden="true">
      <button class="hist-btn" aria-label="Back">&#9664;</button>
      <button class="hist-btn" aria-label="Forward">&#9654;</button>
    </div>

    <div class="brand">
      <span class="brand-title">Zimbra Administration</span>
      <span class="brand-sub">A&nbsp;SYNACOR&nbsp;PRODUCT</span>
    </div>

    <div class="search">
      <button class="search-type" aria-label="Search type">
        <span class="lock">&#128274;</span><span class="caret">&#9662;</span>
      </button>
      <input type="text" aria-label="Search" />
      <button class="search-go" aria-label="Search">&#128269;</button>
    </div>

    <div class="masthead-actions">
      <button class="mh-icon" aria-label="Refresh">&#8635;</button>
      <span class="sep"></span>
      <div class="mh-user" @click="menuOpen = !menuOpen">
        <span>{{ auth.username || 'admin' }}</span>
        <span class="caret">&#9662;</span>
        <ul v-if="menuOpen" class="mh-menu" @click.stop>
          <li @click="logout">Sign Out</li>
        </ul>
      </div>
      <span class="sep"></span>
      <div class="mh-help"><span>Help</span><span class="caret">&#9662;</span></div>
    </div>

    <div class="queue" aria-hidden="true">
      <span class="q-row"><i class="q-dot q-blue"></i>0</span>
      <span class="q-row"><i class="q-dot q-amber"></i>0</span>
      <span class="q-row"><i class="q-dot q-green"></i>0</span>
    </div>
  </header>
</template>

<style scoped>
.masthead {
  display: flex;
  align-items: center;
  height: 44px;
  padding: 0 8px;
  color: var(--masthead-txt);
  background: linear-gradient(to bottom, var(--masthead-top), var(--masthead-bottom));
  border-bottom: 1px solid var(--masthead-border);
  gap: 12px;
}
.hist {
  display: flex;
  gap: 2px;
}
.hist-btn {
  width: 20px;
  height: 20px;
  border: 1px solid #b7c3cc;
  border-radius: 50%;
  background: linear-gradient(to bottom, #ffffff, #e2e8ec);
  color: #3f7cae;
  font-size: 9px;
  line-height: 1;
  cursor: pointer;
}
.brand {
  display: flex;
  flex-direction: column;
  line-height: 1.05;
  min-width: 210px;
}
.brand-title {
  font-size: 19px;
  font-weight: 400;
  letter-spacing: 0.2px;
  color: var(--masthead-txt);
}
.brand-sub {
  font-size: 8px;
  letter-spacing: 1px;
  color: var(--masthead-sub);
}
.search {
  display: flex;
  align-items: center;
  width: 400px;
  height: 24px;
  margin: 0 auto;
  background: #fff;
  border: 1px solid #b9c4cc;
  border-radius: var(--radius);
  padding: 0 2px 0 0;
}
.search-type {
  display: flex;
  align-items: center;
  gap: 2px;
  height: 22px;
  padding: 0 5px;
  border: none;
  border-right: 1px solid #d2d2d2;
  background: linear-gradient(to bottom, #fbfbfb, #ececec);
  cursor: pointer;
  font-size: 10px;
  color: var(--masthead-icon);
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
  font-size: 12px;
  color: var(--masthead-icon);
}
.masthead-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--fs);
  color: var(--masthead-txt);
}
.mh-icon {
  background: none;
  border: none;
  color: var(--masthead-icon);
  font-size: 15px;
  cursor: pointer;
}
.sep {
  width: 1px;
  height: 20px;
  background: #bcc7cf;
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
.queue {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1px;
  padding-left: 6px;
  margin-left: 2px;
  border-left: 1px solid #bcc7cf;
  font-size: 9px;
  color: var(--masthead-sub);
}
.q-row {
  display: flex;
  align-items: center;
  gap: 3px;
}
.q-dot {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  display: inline-block;
}
.q-blue {
  background: #4a90c2;
}
.q-amber {
  background: #e0a83e;
}
.q-green {
  background: #5aa845;
}
</style>
