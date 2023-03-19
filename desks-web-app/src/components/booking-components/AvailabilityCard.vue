<template>

  <v-card :v-if="startOfDay"
          :disabled="isDisabled() || isLoading()"
          :color="calculateColor()"
          :variant="isDisabled() || isLoading() ?'flat':'tonal'"
    >
    <template v-slot:append>
      <v-btn icon="mdi-plus" elevation="0" @click="$emit('addBooking')"></v-btn>
    </template>
    <template v-slot:title>
      {{ startOfDay.format("dddd") }}
    </template>
    <template v-slot:subtitle>
      {{ $format.date(startOfDay)}}
    </template>


    <v-card-text>

      <v-menu
          open-on-hover
          :disabled="isLoading() || this.bookingNumber <= 0"
      >
        <template v-slot:activator="{ props }">
          <v-chip v-if="isLoading()" class="w-75"></v-chip>
          <v-chip v-else v-bind="props">{{ this.bookingNumber }} / {{ this.room.capacity }} desks booked</v-chip>
        </template>

        <booking-list :bookings="getBookingsByRoomAndDay(room.id, startOfDay)" ></booking-list>
      </v-menu>

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
import {th} from "vuetify/locale";
import BookingList from "@/components/booking-components/BookingList.vue";

interface AvailabilityCardState {
  bookingNumber: number
}

export default defineComponent({
  name: 'availability-card',
  components: {BookingList},
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
    ...mapGetters("bookings", ["getBookingsByRoomAndDay", "isLoadingBookingsByRoomAndDay"]),
    ...mapState("bookings", ["bookings"]),
  },
  watch: {

    bookings() {
      this.check()
    },
    room() {
      this.check()
    }
  },
  methods: {

    isDisabled(): boolean {
      return  this.startOfDay.isBefore(moment().startOf("day"))
    },
    isLoading (): boolean {
      return this.isLoadingBookingsByRoomAndDay(this.room.id, this.startOfDay)
    },


    check() {
      const endOfDay = moment(this.startOfDay).endOf("day")
      const overlaps = BookingUtils.findOverlaps(this.getBookingsByRoomAndDay(this.room.id,this.startOfDay ), this.startOfDay, endOfDay)
      this.bookingNumber = new Set(overlaps.map(booking => booking.user.id)).size
    },
    calculateColor(): string {
      if (this.isDisabled() || this.isLoading()) {
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