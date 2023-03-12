<template>

  <v-card :v-if="startOfDay" :disabled="isDisabled(startOfDay)" :color="calculateType(startOfDay)"
          :variant="isDisabled(startOfDay)?'flat':'tonal'">
    <template v-slot:append>
      <v-btn icon="mdi-plus" elevation="0" @click="$emit('addBooking')"></v-btn>
    </template>
    <template v-slot:title>
      {{ startOfDay.format("dddd") }}
    </template>
    <template v-slot:subtitle>
      {{ startOfDay.format("L") }}
    </template>


    <v-card-text>

      <v-chip>{{ this.bookingNumber }} / {{ this.room.capacity }} desks booked</v-chip>

    </v-card-text>
  </v-card>
</template>

<script lang="ts">

import {defineComponent, PropType} from "vue";
import moment from "moment";
import {mapActions, mapGetters, mapState} from "vuex";
import {Moment} from "moment/moment";
import Room from "@/models/Room";
import BookingUtils from "@/utils/booking-utils";

interface AvailabilityCardState {
  bookingNumber: number
}

export default defineComponent({
  name: 'availability-card',
  components: {},
  emits: ["addBooking"],
  props: {
    startOfDay: {
      type: Object as PropType<Moment>,
      required: true
    },
    room: {
      type: Object as PropType<Room>,
      required: true
    }
  },
  data: (): AvailabilityCardState => ({
    bookingNumber: 0
  }),
  computed: {

    ...mapState("defaults", ["bookingDefaults", "bookingDefaultsLoading"]),
    ...mapGetters("bookings", ["getBookingsByRoomAndDay", "isLoadingBookingsByRoomAndDay"]),
    ...mapState("bookings", ["bookings"]),
    ...mapState('user', ['user']),
  },

  watch: {
    async startOfDay() {
      await this.loadAndCheck()
    },
    async room() {
      await this.loadAndCheck()
    },
    bookings() {
      this.check()
    }
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsByRoomAndDate"]),
    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    async loadAndCheck() {

      await this.fetchBookingsByRoomAndDate({roomId: this.room.id, date: this.startOfDay})
    },

    check() {
      const endOfDay = moment(this.startOfDay).endOf("day")
      const overlaps = BookingUtils.findOverlaps(this.bookings, this.startOfDay, endOfDay)
      this.bookingNumber = new Set(overlaps.map(booking => booking.user.id)).size
    },
    calculateType(startOfDay: Moment): string {
      if (this.isDisabled(startOfDay)) {
        return ""
      }
      if (this.bookingNumber >= this.room.capacity) {
        return "red"
      }
      if (this.bookingNumber > this.room.capacity * 0.8) {
        return "orange"
      }

      return "success"

    }
  }
});
</script>

<style scoped>

</style>