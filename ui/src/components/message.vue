<template>
    <div style="padding: 12px; border: 1px solid black; margin-bottom: 12px; white-space: pre;">
        {{ msg.content }}
    </div>
</template>

<script setup lang="ts">
import { Message } from '../api/api';

defineProps<{msg: Message}>()
</script>