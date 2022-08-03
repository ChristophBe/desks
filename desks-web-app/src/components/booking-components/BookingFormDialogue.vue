<template>
  <div v-if="roomsLoading || configurationLoading || defaultsLoading">
    <v-progress-circular
        indeterminate
    ></v-progress-circular>
  </div>
  <booking-form v-else :booking="booking" :defaults="defaults" @close="$emit('close')"></booking-form>
</template>

<script>
import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import BookingForm from "@/components/booking-components/BookingForm";

export default defineComponent( {
  name: "booking-form-dialogue",
  components: {BookingForm},
  props: {
    booking: {
      type: Object
    }
  },
  mounted() {
    this.fetchRooms()
    this.fetchConfiguration()
    this.fetchBookingDefaults()
  },
  computed: {
    ...mapState('rooms', {rooms: 'rooms', roomsLoading: 'loading'}),
    ...mapState('configuration', {configuration: 'configuration', configurationLoading: 'loading'}),
    ...mapState('defaults', {defaults: 'bookingDefaults', defaultsLoading: 'bookingDefaultsLoading'}),
    ...mapGetters('bookings', {getOverlaps: 'getOverlappingBookings'})
  },
  methods: {
    ...mapActions("rooms", ["fetchRooms"]),
    ...mapActions("configuration", ["fetchConfiguration"]),
    ...mapActions("defaults", ["fetchBookingDefaults"])
  }
})
</script>

<style scoped>

</style>
