<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useMailStore } from '../stores/mail'

const auth = useAuthStore()
const mail = useMailStore()

async function onSubmit() {
  const ok = await auth.login()
  if (ok) await mail.loadMailbox()
}
</script>

<template>
  <div
    class="min-h-full flex flex-col items-center pt-24 px-6 pb-10"
    style="background-color:#48525c"
  >
    <form
      class="w-full max-w-md bg-accent border border-accent-bar text-white login-shadow rounded-sm"
      @submit.prevent="onSubmit"
    >
      <div class="flex items-center gap-3 px-5 py-3 bg-accent-bar border-b border-accent-2">
        <div class="text-lg font-bold tracking-tight">
          go-snappymail
          <span v-if="auth.appVersion" class="ml-2 text-xs font-normal opacity-70 font-mono">
            {{ auth.appVersion }}
          </span>
        </div>
      </div>

      <div class="px-6 py-5 flex flex-col gap-3">
        <p v-if="auth.loginErr" class="text-sm text-red-200 bg-red-950/40 border border-red-300/40 px-3 py-2">
          {{ auth.loginErr }}
        </p>

        <label class="text-sm opacity-90">
          Email
          <input
            v-model="auth.loginUser"
            type="email"
            autocomplete="username"
            class="login-input mt-1 block"
            :disabled="auth.loginBusy"
            placeholder="user@test.local"
          />
        </label>

        <label class="text-sm opacity-90">
          Password
          <input
            v-model="auth.loginPwd"
            type="password"
            autocomplete="current-password"
            class="login-input mt-1 block"
            :disabled="auth.loginBusy"
          />
        </label>

        <button
          type="submit"
          class="mt-2 h-9 bg-[#f5f7fa] text-accent-bar font-semibold disabled:opacity-60"
          :disabled="auth.loginBusy"
        >
          {{ auth.loginBusy ? 'Signing in…' : 'Sign in' }}
        </button>
      </div>
    </form>
  </div>
</template>
