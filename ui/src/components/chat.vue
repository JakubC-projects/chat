
<template>
  <div style="height: 100%; display:flex; flex-direction: column; overflow: hidden">
    <h3>Chat {{ chatUid }}</h3>

    <div style="flex-grow: 1; padding: 8px; overflow:scroll; display:flex; flex-direction: column-reverse; justify-content: end;">
        <MessageComp v-for="msg of messages" :msg="msg" :key="msg.id" />
        <!-- <div style="overflow-anchor: auto; height: 70px; background-color: red;"></div> -->
    </div>
    
    <form ref="sendForm" action="GET" @submit="onFormSubmit" style="padding: 12px">
      <input hidden name="type" value="text"/>
      <input hidden name="chat_uid" :value="chatUid"/>
      <div style="display: flex; gap: 12px; align-items:start">
        <textarea ref="input" name="content" style="font-size: 24px; flex-grow: 1; field-sizing: content"></textarea>
        <input type="submit" value="S" style="flex-shrink: 0;"/>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { api, ChatEventListener, type Message } from '../api/api';
import MessageComp from './message.vue'

const props = defineProps<{chatUid: string}>()

const sendForm = ref<HTMLFormElement | null>(null);
const messages = ref<Message[]>([]);
const input = ref<HTMLInputElement | null>(null);

let eventListener: ChatEventListener;
onMounted(() => {
  eventListener = api.listenForMessages(props.chatUid, (m: Message) => {
    messages.value.unshift(m)
  })
})

onUnmounted(() => {
  if(eventListener) {
    eventListener.close()
  }
})

async function onFormSubmit(ev: Event) {
  ev.preventDefault()
  const data = new FormData(ev.target as HTMLFormElement);

  (ev.target as HTMLFormElement).reset()
  input.value?.focus();
  await api.sendMessage(data)
}
</script>

<style>
html {
  height: 100%;
}
body {
  margin: 0;
  height: 100%;
}
</style>