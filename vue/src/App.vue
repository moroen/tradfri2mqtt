<template>
  <q-layout view="hHh LpR fFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <q-avatar>
            <img
              src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg"
            /> </q-avatar
          >Tradfri2MQTT
        </q-toolbar-title>
        <save-button />
      </q-toolbar>
    </q-header>

    <q-drawer show-if-above v-model="leftDrawerOpen" side="left">
      <navigation></navigation>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, onBeforeMount } from "vue";
import Navigation from "./components/Navigation.vue";
import SaveButton from "./components/SaveButton.vue";
import { useStore } from "vuex";

const leftDrawerOpen = ref(false);
const store = useStore();
const footBar = ref(null);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

onBeforeMount(() => {
  console.log("BeforeMount");
  // store.dispatch("getSettings");
  // store.dispatch("startWebSocket");
});

const status = store.state.status;
</script>
