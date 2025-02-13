import { notFound } from 'next/navigation';
import { products } from '@/data/products';
import AddToCartButton from '@/components/AddToCartButton';
import Nav from '@/components/Nav';
import Image from 'next/image';
import Cart from '@/components/Cart';

type ProductPageProps = {
  params: {
    id: string;
  };
};

export default function ProductPage({ params }: ProductPageProps) {
  const product = products.find((p) => p.id === params.id);

  if (!product) {
    notFound();
  }

  return (
    <Nav category={product.category} product={product.name}>
      <main className="max-w-6xl mx-auto py-8 px-4">
        <div className="grid md:grid-cols-2 gap-8">
          <div>
            <Image
              src={product.image}
              alt={product.name}
              width={500}
              height={500}
              className="w-full rounded-lg shadow-md"
            />
          </div>
          <div>
            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <p className="text-gray-600 text-xl mb-6">${product.price}</p>
            <p className="text-gray-700 mb-8">{product.description}</p>
            <AddToCartButton product={product} />
          </div>
        </div>
      </main>
      <Cart />
    </Nav>
  );
}