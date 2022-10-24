<template>

  <v-card class="card" v-if="booking">

    <v-card-title class="d-flex">
      <div class="mr-auto my-auto">Booking Details</div>
      <v-btn icon elevation="0">
        <v-icon>mdi-dots-vertical</v-icon>
        <booking-context-menu activator="parent" :booking="booking" @edit="$emit('edit')"
                              @deleted="$emit('close')"></booking-context-menu>
      </v-btn>
    </v-card-title>
    <v-list>
      <v-list-item>
        <v-list-item-avatar start icon="mdi-calendar"></v-list-item-avatar>
        <v-list-item-header>
          <v-list-item-title>{{ $format.date(booking.start) }}</v-list-item-title>
          <v-list-item-subtitle>{{ $format.timeRange(booking.start, booking.end) }}</v-list-item-subtitle>
        </v-list-item-header>

      </v-list-item>
      <v-list-item>
        <v-list-item-avatar start icon="mdi-map-marker"></v-list-item-avatar>
        <v-list-item-header>

          <v-list-item-title>{{ booking.room.name }}</v-list-item-title>
        </v-list-item-header>
      </v-list-item>
    </v-list>

    <v-card-title class="mt-4">
      Colleagues
      <v-fade-transition>
        <v-progress-circular v-if="!booking || isLoadingBookingsByRoomAndDay(booking.roomId, booking.date)" class=" mr-4" indeterminate size="small"></v-progress-circular>

      </v-fade-transition>
    </v-card-title>
    <also-in-the-room v-if="booking" :date="booking.start" :room-id="booking.room.id"></also-in-the-room>
    <v-card-actions>
      <v-btn @click="$emit('close')">
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import AlsoInTheRoom from "@/components/booking-components/AlsoInTheRoom.vue";
import BookingContextMenu from "@/components/booking-components/BookingContextMenu.vue";
import {mapGetters} from "vuex";

export default defineComponent({
  name: "booking-details",
  components: {AlsoInTheRoom, BookingContextMenu},
  props: {
    booking: {
      type: Object
    }
  },
  computed:{
    ...mapGetters('bookings', ['isLoadingBookingsByRoomAndDay'])
  }
});
</script>

<style scoped>
.card {
  min-width: 400px;
}
</style>
