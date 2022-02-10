<template>
  <div class="q-pa-md" style="max-width: 300px">About</div>

  <q-form>
    <q-input v-model="gateway"></q-input>
  </q-form>
</template>

<script setup>
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";

const store = useStore();
</script>

<script>
export default {
  computed: {
    ...mapFields(["tradfri.gateway", "tradfri.ident"]),
  },
};
</script>
