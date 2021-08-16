<template>
  <q-page padding>
    <!-- content -->
    <div> {{ gateway }} </div>
    <div> {{ id }} </div>
    <div class="q-pa-md">
      <q-btn label="Prompt" color="primary" @click="prompt"></q-btn>
    </div>

  </q-page>
  <psk ref="foo"></psk>
</template>

<script>
import { computed } from "vue";
import { useStore } from "vuex";
import getPSK from "../components/GetPSK.vue";
import { useQuasar } from "quasar";
import { ref } from "vue";

export default {
  components: {
    psk: getPSK
  },

  methods: {
    prompt() {
      this.$q
        .dialog({
          component: getPSK,

          // props forwarded to your custom component
          componentProps: {
            text: "something"
            // ...more..props...
          }
        })
        .onOk(() => {
          console.log("OK");
        })
        .onCancel(() => {
          console.log("Cancel");
        })
        .onDismiss(() => {
          console.log("Called on OK or Cancel");
        });
    }
  },
  mounted() {},

  setup() {
    const $q = useQuasar();
    const $store = useStore();

    const id = computed({
      get: () => $store.state.settings.tradfri.id,
      set: (val) => {
        $store.commit("settings/updateTradfri", {
          section: "tradfri",
          key: "id",
          value: ref(val)
        });
      }
    });

    const gateway = computed({
      get: () => $store.state.settings.tradfri.gateway,
      set: (val) => {
        $store.commit("settings/updateTradfri", {
          section: "tradfri",
          key: "gateway",
          value: ref(val)
        });
      }
    });

    return {
      gateway,
      id
    };
  }
};
</script>
