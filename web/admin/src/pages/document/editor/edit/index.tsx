import { getApiV1NodeDetail } from '@/request/Node';
import { V1NodeDetailResp } from '@/request/types';
import { useAppSelector } from '@/store';
import { Box, Typography, Alert } from '@mui/material';
import { useEffect, useState } from 'react';
import { useOutletContext, useParams } from 'react-router-dom';

import { WrapContext } from '..';
import LoadingEditorWrap from './Loading';
import EditorWrap from './Wrap';

const Edit = () => {
  const { id = '' } = useParams();
  const { kb_id = '' } = useAppSelector(state => state.config);
  const { setNodeDetail } = useOutletContext<WrapContext>();
  const [loading, setLoading] = useState(false);
  const [detail, setDetail] = useState<V1NodeDetailResp | null>(null);
  const [error, setError] = useState<string | null>(null);

  const getDetail = () => {
    setLoading(true);
    setError(null);
    getApiV1NodeDetail({
      id,
      kb_id,
    })
      .then(res => {
        setDetail(res);
        setNodeDetail(res);
        setTimeout(() => {
          window.scrollTo({ top: 0, behavior: 'smooth' });
        }, 0);
      })
      .catch(err => {
        console.error('获取文档详情失败:', err);
        setError('无法找到文档。文档可能已被删除或您没有权限访问。');
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    if (id && kb_id) {
      getDetail();
    }
  }, [id, kb_id]);

  return (
    <Box
      sx={{
        position: 'relative',
        flexGrow: 1,
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        /* Give a remote user a caret */
        '& .collaboration-carets__caret': {
          borderLeft: '1px solid #fff',
          borderRight: '1px solid #fff',
          marginLeft: '-1px',
          marginRight: '-1px',
          pointerEvents: 'none',
          position: 'relative',
          wordBreak: 'normal',
        },
        /* Render the username above the caret */
        '& .collaboration-carets__label': {
          borderRadius: '0 3px 3px 3px',
          color: '#fff',
          fontSize: '12px',
          fontStyle: 'normal',
          fontWeight: '600',
          left: '-1px',
          lineHeight: 'normal',
          padding: '0.1rem 0.3rem',
          position: 'absolute',
          top: '1.4em',
          userSelect: 'none',
          whiteSpace: 'nowrap',
        },
      }}
    >
      {loading ? (
        <LoadingEditorWrap />
      ) : error ? (
        <Box
          sx={{
            flexGrow: 1,
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            p: 4,
          }}
        >
          <Alert severity='error' sx={{ mb: 2, maxWidth: '600px' }}>
            {error}
          </Alert>
          <Typography variant='body1' color='text.secondary'>
            文档ID: {id}
          </Typography>
        </Box>
      ) : (
        detail && <EditorWrap detail={detail} />
      )}
    </Box>
  );
};

export default Edit;
