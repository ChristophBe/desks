<template>
  <div>
    <div v-if="bookingsOfColleagues === undefined || isLoadingBookingsByRoomAndDay(roomId, date)">
      <v-progress-circular class="mx-4 mr-4" indeterminate size="small"></v-progress-circular>
      load bookings of your colleagues
    </div>
    <v-alert v-else-if="bookingsOfColleagues.length <= 0" class="mx-4">There are no other bookings for the same
      office.
    </v-alert>

    <booking-list v-else :bookings="bookingsOfColleagues"></booking-list>
  </div>

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
