<template>
  <q-page padding>
    <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
      <q-input filled v-model="gateway" label="Gateway IP *" hint="Hostname or IP-address of the IKEA-tradfri gateway" lazy-rules :rules="[ val => val && val.length > 0 || 'Please type something']" />

      <div class="row items-center content-stretch">
        <div class="col-8">
          <q-input filled v-model="id" label="Ident *" hint="Ident" lazy-rules :rules="[ val => val && val.length > 0 || 'Please type something']" />
          <q-input filled v-model="key" label="Key *" hint="PSK" lazy-rules :rules="[ val => val && val.length > 0 || 'Please type something']" />
        </div>
        <div class="q-pa-md">
          <q-btn label="Generate ID/PSK" color="primary" @click="promptPSK"></q-btn>
        </div>
      </div>

      <div>
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
      </div>
    </q-form>
  </q-page>
</template>

<script>
import { useQuasar } from "quasar";
import { ref } from "vue";
import { computed } from "vue";
import { useStore } from "vuex";
import getPSK from "../components/GetPSK.vue";

export default {
  data() {
    return {
      testme: "abc"
    };
  },
  methods: {
    promptPSK() {
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
  mounted() {
    this.promptPSK();
  },

  setup() {
    const $store = useStore();

    const $q = useQuasar();

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

    const key = computed({
      get: () => $store.state.settings.tradfri.key,
      set: (val) => {
        $store.commit("settings/updateTradfri", {
          section: "tradfri",
          key: "key",
          value: ref(val)
        });
      }
    });

    return {
      gateway,
      id,
      key,

      onSubmit() {
        if (accept.value !== true) {
          $q.notify({
            color: "red-5",
            textColor: "white",
            icon: "warning",
            message: "You need to accept the license and terms first"
          });
        } else {
          $q.notify({
            color: "green-4",
            textColor: "white",
            icon: "cloud_done",
            message: "Submitted"
          });
        }
      },

      onReset() {
        name.value = null;
        age.value = null;
        accept.value = false;
      }
    };
  }
};
</script>
