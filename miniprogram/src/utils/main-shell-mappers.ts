import type { BackendOrderRecord, BackendTicketRecord, OrderSummary, TicketSummary } from '@/types/main-shell';
import { maskSensitive } from '@/utils/sensitive';

function formatDateLabel(value: string): string {
  const [year, month, day] = value.split('-');
  if (!year || !month || !day) return value;
  return `${month}-${day}`;
}

export function mapBackendOrder(record: BackendOrderRecord): OrderSummary {
  return {
    id: record.id,
    title: record.product_name,
    orderNo: record.order_no,
    status: record.status,
    source: record.source_type,
    travelDateLabel: formatDateLabel(record.travel_date),
    quantityLabel: record.quantity_text,
    amountLabel: record.paid_amount_text,
    hint: record.usage_hint,
    ticketId: record.ticket_id
  };
}

export function mapBackendTicket(record: BackendTicketRecord): TicketSummary {
  return {
    id: record.id,
    status: record.status,
    title: record.title,
    productTitle: record.product_name,
    validDateLabel: formatDateLabel(record.valid_date),
    validTimeLabel: record.valid_time,
    quantityLabel: record.quantity_text,
    maskedCode: maskSensitive(record.ticket_code) as string,
    verifyLocation: record.verify_location,
    notice: record.notice,
    orderNo: record.order_no,
    amountLabel: record.paid_amount_text,
    source: record.source_type
  };
}
