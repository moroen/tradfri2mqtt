<template>
  <q-page padding>
    <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
      narrow-indicator>
      <q-tab name="settings" label="Settings" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="settings">
        <q-form>
          <q-checkbox v-model="enable" label="Enable Interface" />

          <div class="q-pa-md"></div>

          <div class="row justify-start">

            <div class="q-pl-sm">
              <q-input standard v-model.number="port" label="Port *" hint="Interface port" :rules="portRule"
                :disable="!enable" type="number" />
            </div>
          </div>
        </q-form>
      </q-tab-panel>

      <q-tab-panel name="devices">
        <div class="text-h6">Devices</div>
        Lorem ipsum dolor sit amet consectetur adipisicing elit.
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { useQuasar } from "quasar";
import { useStore } from "vuex";
import { computed, ref } from "vue";

import { stringRules, portRule } from "../rules.js";

import getPSK from "../components/GetPSK.vue";
import { mapFields } from "vuex-map-fields";

const store = useStore();
const q = useQuasar();

const tab = ref("settings");

// const nameRules = [val => (val && val.length > 0) || 'Please type something']

function promptPSK() {
  q.dialog({
    component: getPSK,

    // props forwarded to your custom component
    componentProps: {
      text: "something",
      // ...more..props...
    },
  })
    .onOk(() => {
      // console.log("OK");
    })
    .onCancel(() => {
      // console.log("Cancel");
    });
}
</script>

<script>
export default {
  computed: {
    ...mapFields([
      "settings.interface.port",
      "settings.interface.enable",
    ]),
  },
};
</script>
