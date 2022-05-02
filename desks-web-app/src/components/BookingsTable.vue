<template>
  <v-table>
    <thead>
    <tr>
      <th class="text-left w-33">
        Room
      </th>
      <th class="text-left">
        Date
      </th>
      <th class="text-left">
        Start
      </th>
      <th class="text-left">
        End
      </th>
    </tr>
    </thead>
    <tbody>

    <tr
        v-for="item in bookings"
        :key="item.id"
    >
      <td>{{ item.room.name }}</td>
      <td>{{ formatDate(item.start) }}</td>
      <td>{{ formatTime(item.start) }}</td>
      <td>{{ formatTime(item.end) }}</td>
    </tr>
    </tbody>
  </v-table>
</template>

<script>
import {defineComponent} from "vue";
import moment from "moment";


export default defineComponent({
  name: "BookingsTable",
  // type inference enabled
  props: {
    bookings: Array
  },

  methods: {
    formatDate(date) {
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date) {
      return moment(date).format("HH:mm")
    }
  }

})
</script>

<style scoped>

</style>