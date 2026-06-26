import type {
  BackendIslandLockRecord,
  BackendIslandSaleRecord,
  BackendIslandVoyageRecord,
  IslandCruiseScenarioId,
  IslandCruiseViewModel,
  IslandPassenger
} from '@/types/island-cruise';
import { mapIslandLock, mapIslandSale, mapIslandVoyage } from '@/utils/island-cruise-mappers';

const departureDate = '2026-07-15';

const backendVoyages: BackendIslandVoyageRecord[] = [
  {
    voyageId: 880101,
    voyageName: '澳门环岛游夜景航班',
    voyageNo: 'MACAU-NIGHT-1930',
    shipName: '珠澳湾游 1 号',
    departureDate,
    departureTime: '19:30',
    upPortId: 2312,
    upPortName: '湾仔码头',
    downPortId: 2332,
    downPortName: '澳门环岛游A',
    voyageState: '售票',
    cabinPriceList: [
      {
        cabinClassId: 'S',
        cabinClassName: '普通舱',
        cabinTypeCode: 'S',
        count: 24,
        fareTypeList: [
          { fareTypeId: 'adult', fareTypeName: '成人票', price: 88, originalPrice: 128, realName: 1 },
          { fareTypeId: 'child', fareTypeName: '儿童票', price: 58, originalPrice: 88, realName: 1 }
        ]
      }
    ]
  },
  {
    voyageId: 880102,
    voyageName: '珠澳夜游观光航班',
    voyageNo: 'ZHU-AO-2030',
    shipName: '海湾夜色 2 号',
    departureDate,
    departureTime: '20:30',
    upPortId: 2312,
    upPortName: '湾仔码头',
    downPortId: 2317,
    downPortName: '珠澳夜游A',
    voyageState: '售票',
    cabinPriceList: [
      {
        cabinClassId: 'S',
        cabinClassName: '普通舱',
        cabinTypeCode: 'S',
        count: 6,
        fareTypeList: [{ fareTypeId: 'adult', fareTypeName: '成人票', price: 108, originalPrice: 138, realName: 1 }]
      }
    ]
  }
];

const passengers: IslandPassenger[] = [
  {
    name: '陈小珠',
    mobile: '13800138000',
    certTypeId: 556,
    certTypeName: '身份证',
    certNo: '440402199001012345',
    fareId: 'adult'
  }
];

const lockRecord: BackendIslandLockRecord = {
  local_order_no: 'HD202607150001',
  orderNo: 'SUP202607150001',
  ticketNo: 'TICKET20260715001',
  lock_expire_at: '2026-07-15T19:15:00+08:00',
  amount: 88
};

const saleRecord: BackendIslandSaleRecord = {
  local_order_no: 'HD202607150001',
  ticketNo: 'TICKET20260715001',
  codeContent: 'HD0105049143806',
  paid_at: '2026-07-15T18:05:00+08:00'
};

export function normalizeIslandCruiseScenario(scenarioId: string): IslandCruiseScenarioId {
  if (
    scenarioId === 'phase3-success' ||
    scenarioId === 'phase3-no-voyage' ||
    scenarioId === 'phase3-failure' ||
    scenarioId === 'phase3-insufficient-stock' ||
    scenarioId === 'phase3-lock-expired' ||
    scenarioId === 'phase3-passenger-invalid' ||
    scenarioId === 'phase3-sale-failed'
  ) {
    return scenarioId;
  }
  return 'phase3-success';
}

export function getIslandCruiseMock(scenarioId: string): IslandCruiseViewModel {
  const normalizedScenarioId = normalizeIslandCruiseScenario(scenarioId);
  const voyages = backendVoyages.map((record) => mapIslandVoyage(record, departureDate));
  const selectedVoyage = voyages[0] ?? mapIslandVoyage(backendVoyages[0]!, departureDate);
  const amount = selectedVoyage.fares[0]?.price ?? 88;
  const saleFailed = normalizedScenarioId === 'phase3-sale-failed';
  const noVoyage = normalizedScenarioId === 'phase3-no-voyage';
  const failure = normalizedScenarioId === 'phase3-failure';
  const insufficientStock = normalizedScenarioId === 'phase3-insufficient-stock';

  const visibleVoyages = noVoyage || failure ? [] : voyages.map((voyage, index) => {
    if (!insufficientStock || index > 0) return voyage;
    return {
      ...voyage,
      fares: voyage.fares.map((fare) => ({ ...fare, stock: 0 }))
    };
  });

  return {
    dataSource: 'mock',
    scenarioId: normalizedScenarioId,
    hero: {
      title: '澳门环岛游',
      summary: '湾仔码头登船，从海上看澳门天际线、城市灯影与珠澳湾区夜色。',
      imageUrl: '/static/phase2/macau-cruise-night-banner-web.jpg'
    },
    ports: {
      up: { id: 2312, name: '湾仔码头' },
      down: { id: 2332, name: '澳门环岛游A' }
    },
    certTypes: [{ id: 556, name: '身份证' }],
    recommended: noVoyage || failure
      ? null
      : {
          date: departureDate,
          label: '夜景航班推荐',
          firstTime: '19:30',
          count: visibleVoyages.length,
          minPriceLabel: '¥88 起'
        },
    voyages: visibleVoyages,
    selectedVoyageId: visibleVoyages[0]?.id ?? null,
    draft: {
      contactName: normalizedScenarioId === 'phase3-passenger-invalid' ? '' : '陈小珠',
      contactMobile: normalizedScenarioId === 'phase3-passenger-invalid' ? '' : '13800138000',
      passengers: normalizedScenarioId === 'phase3-passenger-invalid' ? [] : passengers,
      quantity: normalizedScenarioId === 'phase3-passenger-invalid' ? 0 : 1,
      totalAmount: normalizedScenarioId === 'phase3-passenger-invalid' ? 0 : amount
    },
    lockedOrder: noVoyage || failure ? null : mapIslandLock(lockRecord, amount, normalizedScenarioId === 'phase3-lock-expired'),
    ticket: noVoyage || failure
      ? null
      : mapIslandSale(saleRecord, selectedVoyage, '陈小珠 成人票 x1', saleFailed),
    serviceNotes: [
      '实名购票：每位乘客都需要姓名、手机号、证件类型和证件号。',
      '座位保留：确认订单后会短时间保留当前票源，请在提示时间内完成支付。',
      '提前检票：建议开航前 20 分钟到达湾仔旅游码头。',
      '天气调整：大风、暴雨或临时停航以码头现场通知为准。'
    ],
    errorMessage: failure ? '环岛游班次服务暂时不可用，请稍后重试。' : undefined
  };
}
