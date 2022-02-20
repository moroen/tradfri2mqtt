<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />

        <q-toolbar-title>
          <q-avatar>
            <img
              src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg"
            /> </q-avatar
          >Tradfri
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer show-if-above v-model="leftDrawerOpen" side="left" bordered>
      <Navigation></Navigation>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer elevated class="bg-grey-8 text-white">
      <q-toolbar>
        {{ status }}
        <q-space />
        <SaveButton></SaveButton>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>

<script setup>
import { ref } from "vue";
import Navigation from "./components/Navigation.vue";
import SaveButton from "./components/SaveButton.vue";
import { useStore } from "vuex";

const leftDrawerOpen = ref(false);
const store = useStore();

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

const status = store.state.status;

store.dispatch("getSettings");
store.dispatch("startWebSocket");
</script>
