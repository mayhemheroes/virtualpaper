/*
 * Virtualpaper is a service to manage users paper documents in virtual format.
 * Copyright (C) 2021  Tero Vierimaa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import {useShowController} from "react-admin";
import {Typography, Card, CardContent, Grid} from "@material-ui/core";


export const AdminView = (props) => {
    const {record} = useShowController({...props, resource:"admin", basePath:"/admin", id:"systeminfo" })

    if (!record) return null;
    return (
        <Grid container spacing={3} alignItems="stretch"  flexGrow={1}>
            <Grid item xs={5}>
                <Card>
                    <CardContent>
                        <Typography variant="h4">Server info</Typography>
                        <Typography variant="h6">{record.name} </Typography>
                        <Typography color="textSecondary" >Version: {record.version}, commit: {record.commit} </Typography>
                        <Typography>Go version: {record.go_version} </Typography>
                        <Typography>Uptime: {record.uptime} </Typography>
                    </CardContent>
                </Card>
            </Grid>
            <Grid item xs={6}>
                <Card>
                    <CardContent>
                        <Typography variant="h4">Server installation</Typography>
                        <Typography>{record.imagemagick_version} </Typography>
                        <Typography>Tesseract version: {record.tesseract_version} </Typography>
                        <Typography>Poppler installed: {record.poppler_installed ? 'Yes': 'No'} </Typography>
                    </CardContent>
                </Card>
            </Grid>
        </Grid>
    );
}


export default AdminView;

