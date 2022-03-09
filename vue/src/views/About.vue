<template>
  <q-btn label="Test" @click="test"></q-btn>

  <scroller
    ref="scrollArea"
    :expand="true"
    :bottom-margin="20"
    style="border: 2px solid black; width: 99%"
  >
    <div v-for="(item, index) in messages" :key="index">{{ item }}</div>
  </scroller>
</template>

<script setup>
import { ref, onMounted } from "vue";
import scroller from "../components/Scroller.vue";

// reactive state
const messages = ref([]);

const scrollArea = ref(null);

const createList = () => {
  for (let i = 0; i < 40; i++) {
    messages.value.push("Line " + i);
  }
};

const test = () => {
  messages.value.push("New message");
  scrollArea.value.scrollToBottom();
};

onMounted(() => {
  createList();
});
</script>
