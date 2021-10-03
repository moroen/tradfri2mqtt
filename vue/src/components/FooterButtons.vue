<template>
  <div>
    <div class="q-pa-md">
      <q-btn label="Save Settings" color="primary" @click="saveSettings"></q-btn>
    </div>
  </div>
</template>

<script>
import { useStore } from "vuex";

export default {
  // name: 'ComponentName',
  methods: {
    saveSettings() {
      console.log("Save Settings clicked");

      this.$store.dispatch("settings/save");
    }
  },

  setup() {
    const $store = useStore();
    return {};
  }
};
</script>
