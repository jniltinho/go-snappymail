<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'
import { useSettingsStore } from '../stores/settings'
import type { SkinId } from '../skins/manifest'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

const showPwd = ref(false)
const keepSignedIn = ref(localStorage.getItem('gsn_keep') === '1')

async function onSubmit() {
  localStorage.setItem('gsn_keep', keepSignedIn.value ? '1' : '0')
  const ok = await auth.login()
  if (ok) await mail.loadMailbox()
}

function onSkinChange(e: Event) {
  settings.skin = (e.target as HTMLSelectElement).value as SkinId
}
</script>

<template>
  <div class="min-h-full flex flex-col items-center pt-24 px-6 pb-10 login-page">
    <form class="w-full max-w-md login-card border login-shadow" @submit.prevent="onSubmit">
      <div class="login-header flex items-center gap-3 px-5 py-3 border-b">
        <div class="text-lg font-bold tracking-tight">
          go-snappymail
          <span class="ml-2 text-xs font-normal opacity-70 font-mono">{{ settings.skinLabel(settings.skin) }}</span>
          <span v-if="auth.appVersion" class="ml-1 text-xs font-normal opacity-60 font-mono">{{ auth.appVersion }}</span>
        </div>
      </div>

      <div class="login-heading px-6 pt-4 text-lg font-bold">Login</div>

      <div class="px-6 py-5 flex flex-col gap-3">
        <p v-if="auth.loginErr" class="login-error text-sm px-3 py-2">
          {{ auth.loginErr }}
        </p>

        <label class="text-sm login-row">
          <span class="login-label">Email</span>
          <input
            v-model="auth.loginUser"
            type="email"
            autocomplete="username"
            class="login-input mt-1 block"
            :disabled="auth.loginBusy"
            placeholder="user@test.local"
          />
        </label>

        <label class="text-sm login-row">
          <span class="login-label">Password</span>
          <span class="login-pwd-wrap block relative">
            <input
              v-model="auth.loginPwd"
              :type="showPwd ? 'text' : 'password'"
              autocomplete="current-password"
              class="login-input mt-1 block"
              :disabled="auth.loginBusy"
            />
            <button type="button" class="login-show" tabindex="-1" @click="showPwd = !showPwd">
              {{ showPwd ? 'Hide' : 'Show' }}
            </button>
          </span>
        </label>

        <div class="flex items-center gap-4 mt-2">
          <button type="submit" class="login-btn h-9 font-semibold disabled:opacity-60" :disabled="auth.loginBusy">
            {{ auth.loginBusy ? 'Signing in…' : 'Login' }}
          </button>
          <label class="login-remember text-sm flex items-center gap-1.5">
            <input v-model="keepSignedIn" type="checkbox" />
            Stay signed in
          </label>
        </div>
      </div>

      <div class="login-footer px-6 py-3 flex items-center gap-3">
        <span class="text-sm">Version</span>
        <select class="login-version" :value="settings.skin" @change="onSkinChange">
          <option v-for="s in settings.availableSkins" :key="s.id" :value="s.id">
            {{ s.id === 'zimbra' ? 'Classic' : s.label }}
          </option>
        </select>
      </div>
    </form>
  </div>
</template>
