<template>
  <q-page padding>
    <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
      narrow-indicator>
      <q-tab name="settings" label="Settings" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="settings">
        <q-form>
          <q-checkbox v-model="enable" label="Enable MQTT" />

          <div class="q-pa-md"></div>

          <div class="row justify-start">
            <div class="col-md-3">
              <q-input standard v-model="host" label="Broker IP *" hint="Hostname or IP-address of the MQTT-broker"
                :rules="stringRules" :disable="!enable" />
            </div>

            <div class="q-pl-sm">
              <q-input standard v-model.number="port" label="Port *" hint="MQTT port" :rules="portRule"
                :disable="!enable" type="number" />
            </div>
          </div>
          <div class="q-pa-md"></div>
          <div class="row justify-start">
            <div class="col-4">
              <q-input standard v-model="commandtopic" label="Command topic"
                hint="Subscribed topic for controlling devices" :rules="stringRules" :disable="!enable" />
            </div>
          </div>

          <div class="row justify-start">
            <div class="col-md-4">
              <q-input standard v-model="discoverytopic" label="Discovery topic"
                hint="Topic for sending device configurations" :rules="stringRules" :disable="!enable" />
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
      "settings.mqtt.host",
      "settings.mqtt.port",
      "settings.mqtt.commandtopic",
      "settings.mqtt.enable",
      "settings.mqtt.discoverytopic",
    ]),
  },
};
</script>
