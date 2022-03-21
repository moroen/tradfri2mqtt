<template>
  <q-page padding style="padding-top: 40px">
    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="gateway">
        <gateway></gateway>
      </q-tab-panel>

      <q-tab-panel name="devices">
        <devices></devices>
      </q-tab-panel>
    </q-tab-panels>
    <q-page-sticky expand position="top">
      <q-toolbar class="bg-white text-white">
        <q-tabs
          v-model="tab"
          dense
          class="text-grey"
          active-color="primary"
          indicator-color="primary"
          align="left"
          narrow-indicator
        >
          <q-tab name="devices" label="Devices" />
          <q-tab name="gateway" label="Gateway" />
        </q-tabs>
      </q-toolbar>
    </q-page-sticky>
  </q-page>
</template>

<script setup>
import { ref } from "vue";

import devices from "../components/DeviceList.vue";
import gateway from "../components/Gateway.vue";

const tab = ref("devices");
</script>
