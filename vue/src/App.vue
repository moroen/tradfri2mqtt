<template>
  <q-layout view="hHh LpR fFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <q-avatar> <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg" /> </q-avatar>Tradfri2MQTT
        </q-toolbar-title>
        <save-button />
      </q-toolbar>
    </q-header>

    <q-drawer show-if-above v-model="leftDrawerOpen" side="left" class="column justify-between no-wrap bg-grey-1">
      <navigation></navigation>
      <q-list>
        <q-item class="col-grow"> </q-item>

        <q-item class="text-caption text-grey-5"> {{ status }}</q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref, onBeforeMount, computed } from "vue";
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

var status = computed(() => {
  return store.state.status;
});

const test = computed({
  get: () => {
    return store.state.status;
  },
});
</script>
