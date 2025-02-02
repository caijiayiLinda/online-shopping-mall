export const categories = [
  'Clothing',
  'Tools', 
  'Toys',
  'Beauty',
  'Pets'
];

export interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  image: string;
  category: string;
  quantity?: number;
}

export const products: Product[] = [
  {
    id: '1',
    name: "Women's Plus Pleated Midi Dress",
    image: '/images/dress.jpg',
    price: 34.98,
    description: 'Material: 100% Polyester',
    category: 'Clothing',
  },
  {
    id: '2',
    name: 'VQJTCVLY Cordless Drill',
    image: '/images/drill.jpg',
    price: 35.49,
    description: '21 Voltage & 2 Variable Speeds',
    category: 'Tools',
  },
  {
    id: '3',
    name: 'Hot Wheels Set of 8 Basic Toy Cars & Trucks',
    image: '/images/toy.jpg',
    price: 8.88,
    description: "It's an instant collection with a set of 8 Hot Wheels, including 1 exclusive vehicle!",
    category: 'Toys',
  },
  {
    id: '4',
    name: 'Maybelline Super Stay Teddy Tint, Long Lasting Matte Lip Stain',
    image: '/images/lip.jpg',
    price: 9.97,
    description: "Meet Super Stay Teddy Tint, Maybelline's teddy-soft Lip tint that lasts. Now you can tint Lips in teddy-soft color for a plush, light feel that lasts all day. This no transfer Lipcolor lasts up to 12 hours",
    category: 'Beauty',
  },
  {
    id: '5',
    name: 'Meow Mix Original Choice Dry Cat Food, 16 Pound Bag',
    image: '/images/cat_food.jpg',
    price: 16.98,
    description: 'Contains one (1) 16-pound bag of Meow Mix Original Choice Dry Cat Food, now with a new look',
    category: 'Pets',
  },
];
