<script setup lang="ts">
import type { VoucherTemplateProps } from './VoucherTemplateDefault.vue'

defineProps<VoucherTemplateProps>()
</script>

<template>
  <div class="voucher-wrap">
    <table class="voucher" style="width: 180px; border: 2px solid black; margin: 2px; display: inline-table; border-collapse: collapse; font-family: Helvetica, Arial, sans-serif; font-size: 12px; color: #000; background: #fff;">
      <tbody>
        <tr>
          <td style="text-align: center; font-size: 14px; font-weight: bold; padding: 4px;">{{ hotspotName }}</td>
        </tr>
        <tr>
          <td style="text-align: center; font-size: 14px; font-weight: bold; border-bottom: 1px solid black; padding: 4px;">
            <img v-if="logo" :src="logo" alt="logo" style="height: 30px; border: 0;"><br>
            {{ new Date().toLocaleString('id-ID') }}
          </td>
        </tr>
        <tr>
          <td style="padding: 4px;">
            <table style="text-align: center; width: 170px; font-size: 12px; border-collapse: collapse;">
              <tbody>
                <tr>
                  <td colspan="2">
                    <table style="width: 100%; border-collapse: collapse;">
                      <template v-if="userMode === 'vc'">
                        <tr>
                          <td style="font-size: 12px; padding: 2px;">Kode Voucher</td>
                        </tr>
                        <tr>
                          <td style="width: 100%; border: 1px solid black; font-weight: bold; font-size: 16px; padding: 4px;">{{ username }}</td>
                        </tr>
                      </template>
                      <template v-else>
                        <tr>
                          <td style="width: 50%; padding: 2px;">Username</td>
                          <td style="padding: 2px;">Password</td>
                        </tr>
                        <tr style="font-size: 14px;">
                          <td style="border: 1px solid black; font-weight: bold; padding: 4px;">{{ username }}</td>
                          <td style="border: 1px solid black; font-weight: bold; padding: 4px;">{{ password }}</td>
                        </tr>
                      </template>
                    </table>
                  </td>
                </tr>
                <tr>
                  <td colspan="2" style="border-top: 1px solid black; font-weight: bold; font-size: 16px; padding: 4px;">
                    {{ validity }} {{ timeLimit }} {{ dataLimit }} {{ price }}
                  </td>
                </tr>
                <tr>
                  <td colspan="2" style="font-weight: bold; font-size: 12px; padding: 2px;">
                    Login: http://{{ dnsName }}
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
