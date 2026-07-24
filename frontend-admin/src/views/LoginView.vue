<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { apiError } from '../api/client'
import loginBanner from '../assets/login-banner.png'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const error = ref('')
const busy = ref(false)

async function submit() {
  error.value = ''
  busy.value = true
  try {
    await auth.login(username.value.trim(), password.value)
    router.push({ name: 'overview' })
  } catch (e) {
    error.value = apiError(e)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-center">
    <div class="login-wrap">
      <div class="login-card">
        <div class="login-banner">
          <img :src="loginBanner" alt="Zimbra Admin Console" />
        </div>
        <form class="login-form" @submit.prevent="submit">
          <div class="login-row">
            <label for="u">Username:</label>
            <input id="u" v-model="username" type="text" autocomplete="username" autofocus />
          </div>
          <div class="login-row">
            <label for="p">Password:</label>
            <input id="p" v-model="password" type="password" autocomplete="current-password" />
          </div>
          <p v-if="error" class="login-error">{{ error }}</p>
          <div class="login-actions">
            <button type="submit" :disabled="busy">Login</button>
          </div>
        </form>
      </div>
      <div class="login-reflection"></div>
    </div>
    </div>
    <footer class="login-footer">
      <p>
        Zimbra :: the leader in open source messaging and collaboration ::
        <a href="#">Blog</a> - <a href="#">Wiki</a> - <a href="#">Forums</a>
      </p>
      <p>Copyright &copy; 2005&ndash;2026. All rights reserved.</p>
    </footer>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(to bottom, var(--login-page-top), var(--login-page-bottom));
}

/* Vertically center the card in the space above the footer, like the legacy. */
.login-center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-wrap {
  width: 500px;
}

.login-card {
  background: linear-gradient(to bottom, var(--login-card-top), var(--login-card-bottom));
  /* top right bottom left — logo sits ~30px from the card's left edge (matches
   * the reference), the form stays right-aligned via .login-form padding-right */
  padding: 30px 40px 44px 30px;
  min-height: 258px;
}

.login-banner {
  margin-bottom: 26px;
}
.login-banner img {
  display: block;
  height: 46px;
  width: auto;
}

/* The form group sits left-of-center under the logo (not flush to the card's
 * right edge), matching the reference. */
.login-form {
  padding-right: 42px;
}
.login-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 14px;
  margin-bottom: 10px;
}
.login-row label {
  color: #fff;
  font-size: var(--fs);
  text-align: right;
}
.login-row input {
  width: 232px;
  height: 22px;
  padding: 2px 4px;
  border: 1px solid #7f9db9;
  background: #fff;
  font-family: var(--font);
  font-size: var(--fs);
  color: var(--txt);
}
.login-row input:focus {
  outline: none;
  border-color: var(--accent);
}

.login-error {
  color: #ffe0e0;
  font-size: var(--fs);
  text-align: right;
  margin: 0 0 8px;
}

.login-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;
}
.login-actions button {
  padding: 3px 14px;
  font-family: var(--font);
  font-size: var(--fs);
  color: var(--txt);
  background: linear-gradient(to bottom, var(--login-btn-top), var(--login-btn-bottom));
  border: 1px solid #fff;
  border-radius: var(--radius);
  box-shadow: 0 1px 3px rgba(50, 50, 50, 0.75);
  cursor: pointer;
}
.login-actions button:active {
  background: linear-gradient(to bottom, var(--login-btn-bottom), var(--login-btn-top));
}
.login-actions button:disabled {
  opacity: 0.6;
  cursor: default;
}

/* Glossy Web-2.0 reflection: a bluer band directly under the card that fades
 * smoothly into the page background. */
.login-reflection {
  height: 88px;
  background: linear-gradient(to bottom, rgba(140, 196, 223, 0.9), rgba(220, 235, 244, 0));
}

.login-footer {
  margin-top: auto;
  padding: 30px 0 30px;
  text-align: center;
  color: var(--txt-muted);
  font-size: var(--fs);
  line-height: 1.6;
}
.login-footer a {
  color: var(--txt-muted);
  text-decoration: underline;
}
</style>
