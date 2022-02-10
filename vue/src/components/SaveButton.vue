<template>
  <q-btn v-on:click="onClick" label="Save" />
</template>

<script setup>
import { useStore } from "vuex";

const store = useStore();

const onClick = function () {
  console.log("Click");
  store.dispatch("saveSettings");
};
</script>
