<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'
import { useSettingsStore } from '../stores/settings'

const auth = useAuthStore()
const mail = useMailStore()
const settings = useSettingsStore()

async function onSubmit() {
  const ok = await auth.login()
  if (ok) await mail.loadMailbox()
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
          <input
            v-model="auth.loginPwd"
            type="password"
            autocomplete="current-password"
            class="login-input mt-1 block"
            :disabled="auth.loginBusy"
          />
        </label>

        <button type="submit" class="login-btn mt-2 h-9 font-semibold disabled:opacity-60" :disabled="auth.loginBusy">
          {{ auth.loginBusy ? 'Signing in…' : 'Login' }}
        </button>
      </div>
    </form>
  </div>
</template>
