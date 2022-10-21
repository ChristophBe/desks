<template>
  <v-expand-transition>
    <v-alert v-if=" bookingsOfColleagues  && bookingsOfColleagues.length <= 0" class="mx-4">There are no other bookings for the same
      office.
    </v-alert>

    <booking-list v-else-if="bookingsOfColleagues" :bookings="bookingsOfColleagues"></booking-list>
  </v-expand-transition>

</template>

<script lang="ts">

import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import {mapActions, mapGetters, mapState} from "vuex";

import BookingList from "@/components/booking-components/BookingList.vue";

interface alsoInTheRoomData {
  bookingsOfColleagues?: Array<Booking>
}

export default defineComponent({
  name: "AlsoInTheRoom",
  components: {BookingList},
  data: (): alsoInTheRoomData => ({
    bookingsOfColleagues: undefined,
  }),
  // type inference enabled
  props: {
    roomId: {
      type: Number
    },
    density: String,
    date: String
  },

  async mounted() {
    this.fetchData()
    this.updateBookingsOfColleagues();
  },
  computed: {
    ...mapState('user', ['user']),
    ...mapState('bookings', ['bookings']),
    ...mapGetters('bookings', ['getBookingsByRoomAndDay', 'isLoadingBookingsByRoomAndDay'])
  },
  watch: {
    roomId() {
      this.fetchData()
    },
    date() {
      this.fetchData()
    },
    bookings() {
      this.updateBookingsOfColleagues();
    }
  },
  methods: {
    ...mapActions('bookings', ['fetchBookingsByRoomAndDate']),
    updateBookingsOfColleagues() {
      if (!this.user) {
        return
      }
      this.bookingsOfColleagues = this.getBookingsByRoomAndDay(this.roomId, this.date)
          .filter((booking: Booking) => booking.user.id !== this.user.id)

    },
    fetchData() {
      this.bookingsOfColleagues = undefined
      if (this.roomId && this.date) {
        this.fetchBookingsByRoomAndDate({roomId: this.roomId, date: this.date})
      }
      this.updateBookingsOfColleagues();
    }
  }
})
</script>

<style scoped>

</style>
