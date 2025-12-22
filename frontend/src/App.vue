<template>
  <div class="app">
    <router-view v-slot="{ Component }">
      <template v-if="Component.name !== 'Login'">
        <Layout>
          <component :is="Component" />
        </Layout>
      </template>
      <template v-else>
        <component :is="Component" />
      </template>
    </router-view>
  </div>
</template>

<script setup>
import Layout from './components/Layout.vue'
</script>

<style scoped>
.app {
  width: 100%;
  height: 100vh;
  font-family: Arial, sans-serif;
  overflow: hidden;
}
</style>