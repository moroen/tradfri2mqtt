<template>
  <q-btn label="Test" @click="test(2)"></q-btn>
</template>

<script setup>
import { useStore } from "vuex";
import { onMounted, ref, computed } from "vue";

const store = useStore();

const test = (type) => {
  console.log(store.getters.devices());
};

onMounted(() => {
  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "devices",
      command: "update",
    })
  );
});
</script>
