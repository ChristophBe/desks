<template>

  <v-list-item class=" justify-end align-center" @click='deleteBooking' :disabled="isCurrentOrPastBooking()">

    <v-list-item-avatar  start icon="mdi-delete"></v-list-item-avatar>
    <v-list-item-title>Remove</v-list-item-title>
  </v-list-item>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import moment from "moment";


export default defineComponent({
  props:{
    booking:{
      type: Object
    },
  },
  methods: {
    isCurrentOrPastBooking(){
      return moment(this.booking?.start).isBefore(moment.now())
    },
    deleteBooking() {
      this.$store.dispatch("bookings/deleteBooking", this.booking)
    }
  }
});
</script>

<style scoped>

</style>
