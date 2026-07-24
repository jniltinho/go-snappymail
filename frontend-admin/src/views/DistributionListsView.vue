<script setup lang="ts">
import ListPane from '../components/ListPane.vue'
// Distribution lists have no dedicated table in the PostfixAdmin schema (a
// multi-target alias serves the same role). The list is shown empty to mirror
// ZimbraAdmin's layout; management happens via Aliases.
</script>

<template>
  <ListPane crumb="Distribution Lists">
    <template #menu>
      <li class="disabled">New</li>
      <li class="disabled">Edit</li>
      <li class="disabled">Delete</li>
    </template>
    <table class="grid">
      <thead>
        <tr>
          <th>Distribution List</th>
          <th>Members</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td colspan="3" class="empty">
            No distribution lists. Use multi-target Aliases to forward to a group.
          </td>
        </tr>
      </tbody>
    </table>
  </ListPane>
</template>
