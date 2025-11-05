import { getShareV1NodeDetail } from '@/request/ShareNode';
import { formatMeta } from '@/utils';
import { ResolvingMetadata } from 'next';

export async function generateMetadata(
  { params }: { params: { id: string } },
  parent: ResolvingMetadata,
) {
  const { id } = params;
  let node = {
    name: '无权访问',
    meta: {
      summary: '无权访问',
    },
  };
  try {
    // @ts-ignore
    node = (await getShareV1NodeDetail({ id })) as any;
  } catch (error) {
    console.log(error);
  }

  return await formatMeta(
    { title: node?.name, description: node?.meta?.summary },
    parent,
  );
}
